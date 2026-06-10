# Go Auth Gateway

A lightweight and extensible API Gateway built with Go, designed for microservices architectures.

It provides centralized authentication, request routing, reverse proxy capabilities, and middleware orchestration, allowing backend services to remain isolated while exposing a single entry point to clients.

## Features

- Reverse Proxy for microservices
- JWT Authentication Middleware
- Request forwarding
- Centralized routing
- Middleware pipeline
- Docker support
- Docker Compose environment
- Modular project structure
- Built with Gin

## Architecture

```
                +------------------+
                |      Client      |
                +--------+---------+
                         |
                         |
                +--------v---------+
                |   Go API Gateway |
                |------------------|
                | Authentication   |
                | JWT Validation   |
                | Middlewares      |
                | Reverse Proxy    |
                +--------+---------+
                         |
        +----------------+----------------+
        |                                 |
+-------v-------+                 +-------v-------+
| User Service  |                 | Order Service |
+---------------+                 +---------------+
```

## Project Structure

```
.
├── cmd/
│   └── gateway/
├── internal/
├── pkg/
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

- `cmd/` → Application entrypoint
- `internal/` → Internal business logic
- `pkg/` → Reusable packages
- `Dockerfile` → Container image
- `docker-compose.yml` → Local environment

## Getting Started

### Clone the repository

```bash
git clone https://github.com/your-user/go-auth-gateway.git

cd go-auth-gateway
```

### Install dependencies

```bash
go mod download
```

### Run locally

```bash
go run ./cmd/gateway
```

Or using Docker:

```bash
docker compose up --build
```

## Authentication

The gateway validates incoming JWT tokens before forwarding requests to downstream services.

## Reverse Proxy

Incoming requests are transparently forwarded to configured backend services while preserving headers and context.

## Technologies

- Go
- Gin
- Reverse Proxy
- JWT
- Docker
- Docker Compose
