# Story Teller

Story Teller est une application web de création d'histoires du soir. L'utilisateur choisit un univers, un ton, un décor, un compagnon et une durée. Le frontend transmet ensuite ces choix au backend, qui peut appeler un webhook externe pour générer l'histoire et conserver les dernières histoires dans un fichier JSON.

> **État actuel :** l'interface et le backend fonctionnent. La génération dynamique nécessite une variable `WEBHOOK_URL`. Sans webhook configuré, le backend renvoie volontairement une histoire de secours.

## Fonctionnalités

- page d'accueil animée ;
- six catégories d'histoires ;
- choix de la durée de lecture ;
- envoi des choix au backend ;
- affichage de l'histoire générée ;
- réglage de la taille du texte ;
- sauvegarde des 50 dernières histoires ;
- consultation des histoires sauvegardées ;
- exécution locale avec Docker Compose.

## Technologies

| Partie | Technologies |
| --- | --- |
| Frontend | Vue 3, TypeScript, Vite |
| Interface | Tailwind CSS, DaisyUI, Lucide |
| Animations | CSS, Vue3 Lottie |
| Backend | Go, API HTTP |
| Authentification du webhook | JWT facultatif |
| Stockage actuel | `backend/data/stories.json` |
| Conteneurs | Docker, Docker Compose, Nginx |

## Architecture

```text
Navigateur
   |
   | /api/story et /api/stories
   v
Frontend Vue servi par Nginx
   |
   v
Backend Go
   |-- appelle WEBHOOK_URL si cette variable existe
   |-- utilise fallbackStory() dans le cas contraire
   `-- sauvegarde les histoires dans data/stories.json
```

Le code du backend accepte plusieurs formes de réponse provenant d'un webhook, notamment :

```json
{
  "title": "Le titre de l'histoire",
  "story": "Le texte de l'histoire"
}
```

Le code contient également un traitement de réponses au format Gemini/n8n. Le workflow n8n et sa configuration ne sont cependant pas inclus dans ce dépôt.

## Structure du projet

```text
story-teller/
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── style.css
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
├── docker-compose.yml
└── README.md
```

## Prérequis

Pour la méthode recommandée :

- Windows 10 ou 11, macOS ou Linux ;
- Docker Desktop avec Docker Compose ;
- Git ;
- VS Code, facultatif mais conseillé.

Pour lancer les deux parties sans Docker, il faut également Node.js et Go.

## Installation avec Docker

### 1. Récupérer le projet

```bash
git clone URL_DU_DEPOT_GITHUB
cd story-teller
```

Remplace `URL_DU_DEPOT_GITHUB` par l'adresse de ton dépôt.

### 2. Créer le fichier `.env`

À la racine du projet, crée un fichier nommé `.env` :

```env
SERVER_ADDR=:8080
PORT=8080
```

Pour connecter le webhook prévu par le backend, ajoute ensuite son adresse :

```env
WEBHOOK_URL=https://exemple.com/webhook/story
```

Ne publie jamais un fichier `.env` contenant des secrets. Le fichier `.gitignore` du projet exclut déjà `.env` de Git.

### 3. Construire et démarrer l'application

```bash
docker compose up --build
```

Ouvre ensuite <http://localhost:4173>.

Pour arrêter l'application, utilise `Ctrl+C`, puis éventuellement :

```bash
docker compose down
```

Après une modification du code, reconstruis les images :

```bash
docker compose up --build
```

## Installation sans Docker

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Vite affiche dans le terminal l'adresse locale du frontend.

### Backend

Dans un autre terminal :

```bash
cd backend
go run .
```

Le backend écoute par défaut sur le port `8080`.

## API

### Vérification du serveur

```http
GET /healthz
```

Réponse attendue :

```text
ok
```

### Générer une histoire

```http
POST /story
Content-Type: application/json
```

Exemple de requête :

```json
{
  "title": "Échos de la caverne lumineuse",
  "tone": "Apaisant",
  "setting": "Grotte cachée",
  "companion": "Renard lumineux",
  "length": "8 min reading time"
}
```

### Lire les histoires sauvegardées

```http
GET /stories
```

## Variables d'environnement du backend

| Variable | Obligatoire | Rôle |
| --- | --- | --- |
| `PORT` | Non | Port du serveur, `8080` par défaut |
| `SERVER_ADDR` | Non | Adresse complète d'écoute, par exemple `:8080` |
| `WEBHOOK_URL` | Non | Adresse du webhook de génération |
| `WEBHOOK_TIMEOUT` | Non | Délai maximal de l'appel, `12s` par défaut |
| `WEBHOOK_AUTH_VALUE` | Non | Valeur complète d'un en-tête d'authentification |
| `WEBHOOK_SECRET_HEADER` | Non | Nom de l'en-tête, `Authorization` par défaut |
| `SECRET_KEY` ou `WEBHOOK_SECRET` | Non | Secret utilisé pour produire un JWT |
| `WEBHOOK_SECRET_SCHEME` | Non | Schéma d'authentification, `Bearer` par défaut |
| `WEBHOOK_JWT_SUB` | Non | Sujet du JWT |
| `WEBHOOK_JWT_AUD` | Non | Audience facultative du JWT |
| `WEBHOOK_JWT_ISS` | Non | Émetteur facultatif du JWT |
| `WEBHOOK_TOKEN_KEY` | Non | Nom d'un paramètre secret ajouté à l'URL |
| `WEBHOOK_TOKEN_VALUE` | Non | Valeur de ce paramètre secret |

## Utiliser Git et créer des commits

Un **commit** est une sauvegarde nommée de l'état du projet. Il permet de suivre les modifications et de revenir à une version précédente.

### Première configuration de Git

À effectuer une seule fois sur l'ordinateur :

```bash
git config --global user.name "Ton Nom"
git config --global user.email "ton-adresse@example.com"
```

Utilise de préférence l'adresse associée à ton compte GitHub.

### Préparer un dossier qui n'est pas encore suivi par Git

Dans le terminal ouvert à la racine de `story-teller` :

```bash
git init
git branch -M main
```

### Cycle habituel d'un commit

1. Observer les fichiers modifiés :

   ```bash
   git status
   ```

2. Examiner les modifications :

   ```bash
   git diff
   ```

3. Ajouter seulement les fichiers souhaités :

   ```bash
   git add README.md
   ```

   Pour ajouter toutes les modifications voulues après les avoir vérifiées :

   ```bash
   git add .
   ```

4. Créer le commit :

   ```bash
   git commit -m "docs: améliorer le README"
   ```

5. Vérifier l'historique :

   ```bash
   git log --oneline
   ```

Un bon commit contient une seule modification logique. Exemples de messages :

```text
feat: ajouter le choix de la durée
fix: corriger l'affichage des histoires sauvegardées
docs: expliquer l'installation avec Docker
test: ajouter les tests de fallback
refactor: simplifier la lecture de la réponse du webhook
```

### Publier le projet sur GitHub

Crée d'abord un dépôt vide sur GitHub, sans ajouter automatiquement de README. Dans le terminal :

```bash
git remote add origin URL_DU_DEPOT_GITHUB
git push -u origin main
```

Pour les publications suivantes :

```bash
git push
```

Si `origin` existe déjà, vérifie son adresse avec :

```bash
git remote -v
```

### Fichiers à ne jamais publier

- `.env` ;
- mots de passe et clés API ;
- jetons d'accès ;
- dossiers `node_modules` ;
- données personnelles contenues dans les histoires sauvegardées.

Avant chaque commit, lance toujours `git status` et vérifie la liste des fichiers.

## Tests et vérifications

Le dépôt fourni ne contient actuellement **aucun fichier de test unitaire**. Les commandes de compilation vérifient que le code peut être construit, mais elles ne remplacent pas des tests unitaires.

### Vérifier le frontend actuel

```bash
cd frontend
npm install
npm run build
```

Cette commande effectue notamment une vérification TypeScript et construit la version de production.

### Vérifier le backend actuel

```bash
cd backend
go test ./...
```

Tant qu'aucun fichier `*_test.go` n'existe, Go peut afficher `[no test files]`. Ce message n'est pas une erreur : il indique simplement qu'aucun test n'a encore été écrit.

### Qu'est-ce qu'un test unitaire ?

Un test unitaire vérifie automatiquement une petite fonction isolée :

- on lui donne une valeur d'entrée ;
- on connaît le résultat attendu ;
- le test échoue si le résultat réel est différent.

### Premier exemple de test Go

Crée un fichier `backend/main_test.go` :

```go
package main

import "testing"

func TestFallbackReturnsOriginalValue(t *testing.T) {
	result := fallback("Aventure", "Valeur par défaut")

	if result != "Aventure" {
		t.Fatalf("résultat attendu %q, résultat obtenu %q", "Aventure", result)
	}
}

func TestFallbackReturnsDefaultValue(t *testing.T) {
	result := fallback("", "Valeur par défaut")

	if result != "Valeur par défaut" {
		t.Fatalf("résultat attendu %q, résultat obtenu %q", "Valeur par défaut", result)
	}
}
```

Puis exécute :

```bash
cd backend
go test -v ./...
```

Le résultat doit contenir `PASS`.

### Organisation conseillée pour chaque modification

```text
1. Modifier une petite partie du code
2. Enregistrer les fichiers
3. Lancer les tests
4. Lancer la compilation
5. Vérifier l'application dans le navigateur
6. Examiner git status et git diff
7. Créer le commit
8. Envoyer le commit avec git push
```

Exemple :

```bash
cd backend
go test -v ./...
cd ../frontend
npm run build
cd ..
git status
git diff
git add .
git commit -m "test: ajouter les premiers tests du backend"
git push
```

## Limites actuelles

- le workflow externe de génération n'est pas fourni ;
- sans `WEBHOOK_URL`, une histoire de secours est affichée ;
- les histoires sont stockées dans un fichier JSON et non dans une base de données ;
- le volume de données du backend n'est pas encore déclaré dans `docker-compose.yml` ;
- aucune suite de tests automatisés n'est encore incluse.

## Licence

Aucune licence n'est fournie dans le dépôt actuel. Avant une diffusion ou une réutilisation publique, ajoute un fichier `LICENSE` après avoir vérifié les droits sur le code et les illustrations.
