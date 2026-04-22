# Homeserver Project Overview

A Go-based home server project providing web and file management services.

## Core Technologies
- **Language:** Go 1.23.4
- **Database:** PostgreSQL (using `pgx/v5` driver and connection pooling)
- **Authentication:** JWT-based auth (`golang-jwt/jwt/v5`)
- **Frontend:** HTMX, Vanilla CSS, and SVG icons
- **Environment Management:** `godotenv`
- **Development Tooling:** `air` for live reloading

## Project Structure
- `cmd/`: Entry points for the application.
    - `web/`: Main web server handling user interface and API requests.
    - `file/`: Dedicated file serving and management service.
    - `checks/`: Database and environment health checks.
- `db/`: Bare database layer.
    - Consolidated repository files (e.g., `user_repo.go`, `file_repo.go`).
    - `database.go`: God object for database access (`db.Users`, `db.Files`, etc.).
    - `tables.go`: Database schema definitions.
- `internal/`: Private application code.
    - `server/`: HTTP server initialization and running.
    - `routes/`: Centralized routing logic for both Web and File servers.
    - `setup/`: Application initialization logic.
    - `templates/`: Go template definitions for server-side rendering.
- `pkg/`: Publicly accessible packages.
    - `handlers/`: HTTP request handlers (split by service).
    - `middleware/`: HTTP middleware for auth, logging, etc.
    - `parser/`: Path and request parsing utilities.
    - `utils/`: Flattened utility package for forms, alerts, error responses, etc.
- `static/`: Frontend assets (CSS, JS, SVG, HTML fragments).
- `constants/`: Flattened global constants (`constants.go`).
- `migrations/`: Database schema migrations and scripts.

## Building and Running

### Prerequisites
- Go 1.23+
- PostgreSQL database
- `.env` file in the root directory (see `constants/constants.go` for keys)

### Development
Use `air` for hot reloading:
```bash
# Run Web Server
air -c .air.web.toml

# Run File Server
air -c .air.file.toml
```

### Production
Build the binaries:
```bash
go build -o web_server ./cmd/web/main.go
go build -o file_server ./cmd/file/main.go
```

## Development Conventions
- **Bare Repository Pattern:** Repositories use raw SQL and explicit scanning.
- **God Object:** Access all data via the `db.Database` struct passed to handlers.
- **Validation:** Use `Validate()` methods on entity structs.
- **Flat Structure:** Keep package hierarchies shallow. Avoid single-function packages.
- **HTMX:** Favor HTMX for dynamic frontend updates.
