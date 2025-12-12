# NOTORA - Auth Backend Starter

This is a production-ready starter for NOTORA's authentication backend (Go + Gin + SQLite).

## Quickstart

1. Copy `.env.example` to `.env` and update secrets.
2. Build:
   ```
   go build ./cmd/notora-server
   ```
3. Run:
   ```
   ./notora-server
   ```

Or build the Docker image:
```
docker build -t notora-server .
docker run -p 8080:8080 -v notora-data:/app/data --env-file .env notora-server
```
