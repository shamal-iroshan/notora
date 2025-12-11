# NOTORA  
A simple, fast, self-hosted Markdown note-taking app built with Go, React, and SQLite.

## About
NOTORA is designed to be a lightweight note application that you can easily run on your own server.  
It uses Go for the backend, React (Vite) for the frontend, and SQLite for local storage.

## Tech Stack
- **Go** (API server)
- **React + Vite** (frontend)
- **SQLite** (database)
- **Docker** (optional, for deployment)

## Running Locally

### Backend
```bash
go run main.go
```

### Frontend
```bash
npm install
npm run dev
```

## Docker

Build:
```bash
docker build -t notora .
```

Run:
```bash
docker run -p 8080:8080 notora
```

## License
MIT License. See the LICENSE file for more details.