# 🚀 Go Microservices + API Gateway

This project demonstrates a production-style microservices architecture built with Go, featuring an API Gateway responsible for routing, authentication, and centralized control.

## 🧱 Architecture Overview

The system is composed of:

- API Gateway (reverse proxy + middleware)
- User Service
- Order Service

All services are containerized using Docker and communicate via HTTP.

## 🎯 Goals

- Demonstrate microservices architecture patterns
- Implement API Gateway pattern
- Apply middleware for authentication
- Structure Go services using clean architecture principles

## ⚙️ Tech Stack

- Go (Golang)
- Gin Web Framework
- Docker & Docker Compose
- Reverse Proxy (net/http)

## 🔐 Features

- Centralized routing via API Gateway
- Authentication middleware (JWT-ready)
- Service isolation
- Scalable architecture design

---
