# Multi-stage build for the Go backend
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod ./
# go.sum may be absent; go mod download will generate it if needed
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 8080
CMD ["/app/server"]
