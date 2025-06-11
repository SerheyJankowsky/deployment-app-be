# Deployment App Backend (Go)

## Overview

This project is a backend service for managing deployments, servers, containers, scripts, secrets, and user authentication. It is built with Go, uses the Fiber web framework, GORM for ORM, and supports PostgreSQL as the database. The application is modular, with each core feature implemented as a separate module.

## Features

- **User Authentication**: Register, login, JWT-based authentication, and user management.
- **Secrets Management**: Securely store and manage secrets per user.
- **Servers Management**: CRUD operations for server entities.
- **Containers Management**: CRUD operations for containers associated with servers.
- **Scripts Management**: Manage deployment and utility scripts.
- **Database Migrations**: Automated migrations using Goose and GORM.
- **RESTful API**: All features are exposed via a versioned REST API (`/api/v1`).
- **Dockerized**: Ready for containerized deployment.

## Project Structure

```
.
├── cmd/
│   ├── app/         # Main application entrypoint
│   └── db/          # Database migration entrypoint and migration files
├── modules/         # Feature modules (auth, users, servers, containers, scripts, secrets, etc.)
├── libs/            # Shared libraries (JWT, encryption, etc.)
├── Dockerfile       # Multi-stage Docker build
├── run-app.sh       # Entrypoint script for Docker
├── go.mod           # Go module definition
└── go.sum           # Go dependencies lock file
```

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL database
- Docker (optional, for containerized deployment)

### Environment Variables

Set the following environment variable:

- `DATABASE_URL`: PostgreSQL connection string (e.g., `postgres://user:password@host:port/dbname?sslmode=disable`)

### Local Development

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd deployment-app-be-go
   ```
2. **Install dependencies:**
   ```sh
   go mod download
   ```
3. **Run database migrations:**
   ```sh
   go run cmd/db/main.go
   ```
4. **Start the application:**
   ```sh
   go run cmd/app/main.go
   ```

### Docker Deployment

1. **Build and run with Docker:**
   ```sh
   docker build -t deployment-app .
   docker run -e DATABASE_URL=... -p 8080:8080 deployment-app
   ```

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Auth

- `POST /auth/login` — Login
- `POST /auth/register` — Register
- `POST /auth/refresh` — Refresh JWT
- `GET /auth/me` — Get current user info

### Users

- `PATCH /users/` — Update user
- `DELETE /users/` — Delete user

### Secrets

- `GET /secrets/` — List secrets
- `POST /secrets/` — Create secret
- `PATCH /secrets/:id` — Update secret
- `DELETE /secrets/:id` — Delete secret

### Servers

- `GET /servers/` — List servers
- `POST /servers/` — Create server
- `PATCH /servers/:id` — Update server
- `DELETE /servers/:id` — Delete server

### Containers

- `GET /containers/` — List containers
- `POST /containers/` — Create container
- `PATCH /containers/:id` — Update container
- `DELETE /containers/:id` — Delete container

### Scripts

- `GET /scripts/` — List scripts
- `POST /scripts/` — Create script
- `PATCH /scripts/:id` — Update script
- `DELETE /scripts/:id` — Delete script

### Health Check

- `GET /health` — Returns `OK` if the service is running

## Database Migrations

Migrations are located in `cmd/db/migration/` and are run automatically on startup or can be run manually:

```sh
go run cmd/db/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

## License

[MIT](LICENSE)
