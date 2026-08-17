package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// StoryRequest matches the payload coming from the frontend.
type StoryRequest struct {
	Title     string `json:"title"`
	Tone      string `json:"tone"`
	Setting   string `json:"setting"`
	Companion string `json:"companion"`
	Length    string `json:"length"`
}

// StoryResponse is returned to the frontend.
type StoryResponse struct {
	Title          string `json:"title"`
	Story          string `json:"story"`
	LengthEstimate string `json:"lengthEstimate,omitempty"`
}

// SavedStory represents a persisted story.
type SavedStory struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Story     string    `json:"story"`
	Tone      string    `json:"tone"`
	Length    string    `json:"length"`
	CreatedAt time.Time `json:"createdAt"`
}

// webhookResponse captures the fields we might receive back.
type webhookResponse struct {
	Title   string `json:"title"`
	Story   string `json:"story"`
	Text    string `json:"text"`
	Content string `json:"content"`
}

type httpStatusError struct {
	Code int
	Body string
}

func (e *httpStatusError) Error() string {
	if e.Body != "" {
		return http.StatusText(e.Code) + ": " + e.Body
	}
	return http.StatusText(e.Code)
}

var (
	storiesFile  = "data/stories.json"
	storiesMutex sync.Mutex
)

func main() {
	loadDotEnv(".env")

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("warning: failed to create data directory: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/story", handleStory)
	mux.HandleFunc("/stories", handleGetStories)

	addr := getEnv("SERVER_ADDR", "")
	if addr == "" {
		port := getEnv("PORT", "8080")
		if strings.HasPrefix(port, ":") {
			addr = port
		} else {
			addr = ":" + port
		}
	}
	log.Printf("story backend listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleGetStories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stories, err := getSavedStories()
	if err != nil {
		log.Printf("failed to get stories: %v", err)
		http.Error(w, "failed to retrieve stories", http.StatusInternalServerError)
		return
	}
	writeJSON(w, stories)
}

func handleStory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), getDurationEnv("WEBHOOK_TIMEOUT", 12*time.Second))
	defer cancel()

	var req StoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	resp := StoryResponse{
		Title:          fallback(req.Title, "Tonight's Story"),
		LengthEstimate: req.Length,
		Story:          "",
	}

	webhookURL := os.Getenv("WEBHOOK_URL")

	if webhookURL == "" {
		log.Printf("no WEBHOOK_URL configured, returning placeholder story")
		resp.Story = fallbackStory(req)
	} else if story, title, err := callWebhook(ctx, webhookURL, req); err == nil {
		resp.Story = story
		if title != "" {
			resp.Title = title
		}
	} else {
		log.Printf("webhook call failed: %v (falling back to placeholder)", err)
		if resp.Story == "" {
			resp.Story = fallbackStory(req)
		}
	}

	// Save the generated story
	if resp.Story != "" && !strings.Contains(resp.Story, "Les tisserands d'histoires dorment encore") {
		saved := SavedStory{
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			Title:     resp.Title,
			Story:     resp.Story,
			Tone:      req.Tone,
			Length:    req.Length,
			CreatedAt: time.Now(),
		}
		go func() {
			if err := saveStory(saved); err != nil {
				log.Printf("failed to save story: %v", err)
			}
		}()
	}

	writeJSON(w, resp)
}

func saveStory(story SavedStory) error {
	storiesMutex.Lock()
	defer storiesMutex.Unlock()

	var stories []SavedStory
	if data, err := os.ReadFile(storiesFile); err == nil {
		_ = json.Unmarshal(data, &stories)
	}

	// Prepend new story
	stories = append([]SavedStory{story}, stories...)

	// Limit to last 50 stories to keep file size manageable
	if len(stories) > 50 {
		stories = stories[:50]
	}

	data, err := json.MarshalIndent(stories, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(storiesFile, data, 0644)
}

func getSavedStories() ([]SavedStory, error) {
	storiesMutex.Lock()
	defer storiesMutex.Unlock()

	var stories []SavedStory
	data, err := os.ReadFile(storiesFile)
	if os.IsNotExist(err) {
		return []SavedStory{}, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &stories); err != nil {
		return nil, err
	}
	return stories, nil
}

func callWebhook(ctx context.Context, webhookURL string, req StoryRequest) (story string, title string, err error) {
	httpReq, err := buildWebhookRequest(ctx, webhookURL, req)
	if err != nil {
		return "", "", err
	}
	httpClient := &http.Client{Timeout: getDurationEnv("WEBHOOK_TIMEOUT", 12*time.Second)}

	res, err := httpClient.Do(httpReq)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		bodyStr := string(raw)
		log.Printf("webhook non-200 status=%d body=%s", res.StatusCode, bodyStr)
		return "", "", &httpStatusError{Code: res.StatusCode, Body: bodyStr}
	}

	// Try JSON formats
	if body, title, ok := parseJSONStory(raw); ok {
		return body, title, nil
	}

	// Fallback: treat raw body as story text (e.g., Markdown or plain)
	plain := strings.TrimSpace(string(raw))
	if plain != "" {
		return plain, headingFromBody(plain), nil
	}

	return "", "", fmt.Errorf("webhook response empty or unparseable")
}

func headingFromBody(body string) string {
	lines := strings.Split(body, "\n")
	if len(lines) == 0 {
		return ""
	}
	first := strings.TrimSpace(lines[0])
	if strings.HasPrefix(first, "#") {
		return strings.TrimSpace(strings.TrimLeft(first, "#"))
	}
	return ""
}

func parseJSONStory(raw []byte) (body string, title string, ok bool) {
	// Simple JSON response
	var webhookRes webhookResponse
	if err := json.Unmarshal(raw, &webhookRes); err == nil {
		body = fallback(webhookRes.Story, fallback(webhookRes.Text, webhookRes.Content))
		title = webhookRes.Title
		if body != "" {
			if title == "" {
				title = headingFromBody(body)
			}
			return body, title, true
		}
	}

	// Parts at root: { parts: [{ text: "..."}] }
	var partsWrapper struct {
		Parts []textPart `json:"parts"`
	}
	if err := json.Unmarshal(raw, &partsWrapper); err == nil && len(partsWrapper.Parts) > 0 {
		body = mergeParts(partsWrapper.Parts)
		if body != "" {
			return body, headingFromBody(body), true
		}
	}

	// Gemini/n8n Array format: [{ "content": { "parts": [...] } }]
	// Also handles variants where parts are at the top or text is in content fields.
	type FlexibleContent struct {
		Parts []textPart `json:"parts"`
		Title string     `json:"title"`
		Story string     `json:"story"`
		Text  string     `json:"text"`
	}
	type FlexibleCandidate struct {
		Content FlexibleContent `json:"content"`
		Parts   []textPart      `json:"parts"`
	}

	var candidates []FlexibleCandidate
	if err := json.Unmarshal(raw, &candidates); err == nil && len(candidates) > 0 {
		for _, item := range candidates {
			// 1. Try content.parts (Standard Gemini)
			if len(item.Content.Parts) > 0 {
				body = mergeParts(item.Content.Parts)
				title = headingFromBody(body)
				return body, title, true
			}
			// 2. Try top-level parts
			if len(item.Parts) > 0 {
				body = mergeParts(item.Parts)
				title = headingFromBody(body)
				return body, title, true
			}
			// 3. Try content fields
			if item.Content.Story != "" || item.Content.Text != "" {
				body = fallback(item.Content.Story, item.Content.Text)
				title = headingFromBody(body)
				return body, title, true
			}
		}
	}

	// Generic fallback: search nested maps/arrays for text fields
	var generic any
	if err := json.Unmarshal(raw, &generic); err == nil {
		texts := collectTexts(generic)
		if len(texts) > 0 {
			body = mergeParts(texts)
			title = headingFromBody(body)
			return body, title, true
		}
	}

	return "", "", false
}

type textPart struct {
	Text string `json:"text"`
}

func mergeParts(parts []textPart) string {
	var buf strings.Builder
	for i, p := range parts {
		segment := strings.TrimSpace(p.Text)
		if segment == "" {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(segment)
		if i == 0 && strings.HasPrefix(segment, "###") {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func collectTexts(v any) []textPart {
	switch t := v.(type) {
	case map[string]any:
		// direct text field
		if val, ok := t["text"]; ok {
			if s, ok2 := val.(string); ok2 {
				return []textPart{{Text: s}}
			}
		}
		// nested maps/arrays
		var acc []textPart
		for _, v2 := range t {
			acc = append(acc, collectTexts(v2)...)
		}
		return acc
	case []any:
		var acc []textPart
		for _, v2 := range t {
			acc = append(acc, collectTexts(v2)...)
		}
		return acc
	default:
		return nil
	}
}

func buildJWT(secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": getEnv("WEBHOOK_JWT_SUB", "story-backend"),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	}

	// Optional: audience & issuer if you want to match n8n config
	if aud := getEnv("WEBHOOK_JWT_AUD", ""); aud != "" {
		claims["aud"] = aud
	}
	if iss := getEnv("WEBHOOK_JWT_ISS", ""); iss != "" {
		claims["iss"] = iss
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func buildWebhookRequest(ctx context.Context, webhookURL string, req StoryRequest) (*http.Request, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Optional query token e.g., ?token=abc
	if key := getEnv("WEBHOOK_TOKEN_KEY", ""); key != "" {
		if val := getEnv("WEBHOOK_TOKEN_VALUE", getEnv("SECRET_KEY", "")); val != "" {
			q := "?"
			if strings.Contains(webhookURL, "?") {
				q = "&"
			}
			webhookURL = webhookURL + q + key + "=" + val
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Custom header value (full)
	if authValue := getEnv("WEBHOOK_AUTH_VALUE", ""); authValue != "" {
		header := getEnv("WEBHOOK_SECRET_HEADER", "Authorization")
		httpReq.Header.Set(header, authValue)
		return httpReq, nil
	}

	// JWT using SECRET_KEY / WEBHOOK_SECRET as signing secret
	if secret := getEnv("SECRET_KEY", getEnv("WEBHOOK_SECRET", "")); secret != "" {
		token, err := buildJWT(secret)
		if err != nil {
			return nil, fmt.Errorf("failed to build JWT: %w", err)
		}

		header := getEnv("WEBHOOK_SECRET_HEADER", "Authorization")
		scheme := getEnv("WEBHOOK_SECRET_SCHEME", "Bearer")

		val := token
		if scheme != "" {
			val = scheme + " " + token
		}

		httpReq.Header.Set(header, val)
	}

	return httpReq, nil
}

func fallbackStory(req StoryRequest) string {
	return "Les tisserands d'histoires dorment encore, voici un petit conte en attendant leur réveil.\n\n" +
		"Titre : " + fallback(req.Title, "Une aventure sans titre") + "\n" +
		"Ton : " + fallback(req.Tone, "Calme") + "\n" +
		"Décor : " + fallback(req.Setting, "Une prairie tranquille") + "\n" +
		"Compagnon : " + fallback(req.Companion, "Une chouette amicale") + "\n\n" +
		"Ferme les yeux, respire doucement et imagine l'histoire qui se déroule pour toi."
}

func fallback(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getDurationEnv(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		if n, convErr := strconv.Atoi(val); convErr == nil {
			return time.Duration(n) * time.Second
		}
		log.Printf("invalid duration for %s=%s, using default %v", key, val, def)
		return def
	}
	return dur
}

// loadDotEnv reads a simple KEY=VALUE file and sets env vars if not already set.
func loadDotEnv(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) != 2 {
			continue
		}
		key := string(bytes.TrimSpace(parts[0]))
		val := string(bytes.TrimSpace(parts[1]))
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}
