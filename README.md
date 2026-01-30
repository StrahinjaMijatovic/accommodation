# Accommodation Booking Platform

Microservices-based accommodation booking platform with role-based access control (Host/Guest).

## Technologies

**Backend:** Go, Gorilla Mux, JWT
**Frontend:** Angular 16, TypeScript
**Databases:** MongoDB, Cassandra, Neo4j, Redis
**DevOps:** Docker, Docker Compose, Traefik

## Architecture

| Service | Database | Description |
|---------|----------|-------------|
| auth-service | MongoDB | Registration, login, JWT authentication |
| profile-service | MongoDB | User profile management |
| accommodation-service | Cassandra | Accommodation CRUD, search, pricing |
| reservation-service | Cassandra | Booking management |
| rating-service | Cassandra | Host and accommodation ratings |
| notification-service | Neo4j | Graph-based notifications |

**Redis** is used for image caching.

## Features

- JWT authentication with role-based access (Host/Guest)
- Accommodation listing with search and filtering
- Dynamic pricing (per-guest / per-unit)
- Reservation management
- Rating system for hosts and accommodations
- Notification system

## Setup

```bash
# Backend
cd Backend
docker-compose up --build

# Frontend
cd Frontend/frontend
npm install
ng serve
```

Application runs at `http://localhost:4200`