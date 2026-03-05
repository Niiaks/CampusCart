# CampusCart -- Server Application

The CampusCart server is a Go-based REST API that powers the CampusCart student marketplace. It handles user authentication with session-based cookies, listing management, brand profiles, hierarchical category systems with dynamic attributes, file uploads via Cloudinary, transactional emails via Resend, background job processing, and comprehensive observability through structured logging and New Relic APM.

---

## Table of Contents

- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Development Tasks](#development-tasks)
- [Architecture](#architecture)
- [Request Lifecycle](#request-lifecycle)
- [Configuration](#configuration)
- [Database](#database)
- [Migrations](#migrations)
- [Seeding](#seeding)
- [Authentication and Sessions](#authentication-and-sessions)
- [API Endpoints](#api-endpoints)
- [Middleware Chain](#middleware-chain)
- [Error Handling](#error-handling)
- [Validation](#validation)
- [File Uploads](#file-uploads)
- [Email System](#email-system)
- [Background Jobs](#background-jobs)
- [Logging](#logging)
- [Observability](#observability)
- [Domain Models](#domain-models)
- [Repository Layer](#repository-layer)
- [Service Layer](#service-layer)
- [Handler Layer](#handler-layer)
- [OpenAPI Documentation](#openapi-documentation)
- [Uploader CLI Tool](#uploader-cli-tool)
- [Code Conventions](#code-conventions)
- [Deployment](#deployment)

---

## Tech Stack

| Technology              | Version | Purpose                                                            |
| ----------------------- | ------- | ------------------------------------------------------------------ |
| Go                      | 1.25    | Server programming language                                        |
| chi                     | v5      | Lightweight, idiomatic HTTP router with composable middleware      |
| pgx                     | v5      | High-performance PostgreSQL driver with native connection pooling  |
| tern                    | v2      | SQL migration tool with embedded file support                      |
| go-redis                | v9      | Redis client for caching, sessions, and job queue backend          |
| Asynq                   | latest  | Distributed background task queue built on Redis                   |
| Cloudinary              | v2      | Cloud-based image and video management (upload, transform, delete) |
| Resend                  | v2      | Transactional email delivery API                                   |
| zerolog                 | latest  | Zero-allocation structured JSON logging                            |
| New Relic Go Agent      | latest  | Application performance monitoring, distributed tracing            |
| go-playground/validator | v10     | Struct tag-based request validation with custom error messages     |
| bcrypt                  | latest  | Password hashing using adaptive cost function                      |
| koanf                   | latest  | Configuration loading from environment variables                   |
| httprate                | latest  | HTTP rate limiting middleware                                      |

## Project Structure

```
apps/server/
  |-- cmd/                              # Application entry points
  |    |-- campusCart/
  |    |    |-- main.go                 # Primary API server entry point
  |    |-- uploader/
  |         |-- main.go                 # CLI file upload utility
  |
  |-- internal/                         # Private application code (not importable externally)
  |    |-- config/
  |    |    |-- config.go               # Configuration struct definitions, env var loading, validation
  |    |    |-- observability.go        # New Relic, logging, and health check configuration
  |    |
  |    |-- database/
  |    |    |-- database.go             # PostgreSQL connection pool initialization and configuration
  |    |    |-- migrator.go             # Embedded SQL migration runner using tern
  |    |    |-- pinger.go               # Database health check interface
  |    |    |-- migrations/             # SQL migration files (001 through 010)
  |    |
  |    |-- err/
  |    |    |-- types.go                # HTTPError struct with code, message, status, field errors, redirect
  |    |    |-- http.go                 # Error factory functions (Unauthorized, Forbidden, BadRequest, etc.)
  |    |
  |    |-- handler/
  |    |    |-- base.go                 # Generic handler framework with observability and response serialization
  |    |    |-- handlers.go             # Handler aggregator struct
  |    |    |-- auth.go                 # Authentication handlers (login, register, verify, logout, me)
  |    |    |-- listing.go              # Listing CRUD handlers and upload signature endpoint
  |    |    |-- brand.go                # Brand profile handlers (get, update with file upload)
  |    |    |-- category.go             # Category CRUD handlers (admin-only create/update/delete)
  |    |    |-- health.go               # Health check endpoint (database + Redis)
  |    |    |-- openapi.go              # OpenAPI documentation page handler
  |    |
  |    |-- lib/
  |    |    |-- email/
  |    |    |    |-- client.go          # Resend email client initialization and HTML template rendering
  |    |    |    |-- emails.go          # Email functions (welcome, verification code)
  |    |    |    |-- template.go        # Template name constants
  |    |    |-- file/
  |    |    |    |-- client.go          # Cloudinary client (upload, delete, direct upload signature)
  |    |    |-- job/
  |    |    |    |-- job.go             # Asynq job service (3 priority queues, 10 concurrency)
  |    |    |    |-- email.go           # Email job type definitions and handlers
  |    |    |-- tokenhash/
  |    |         |-- tokenhash.go       # SHA-256 token hashing utility
  |    |
  |    |-- logger/
  |    |    |-- logger.go               # Structured logging setup (console/JSON, New Relic forwarding)
  |    |
  |    |-- middleware/
  |    |    |-- middleware.go            # Middleware aggregator struct
  |    |    |-- auth.go                 # Session authentication with sliding expiration
  |    |    |-- admin.go                # Admin role authorization check
  |    |    |-- global.go               # CORS, request logging, panic recovery, security headers
  |    |    |-- context.go              # Logger context enrichment (request ID, IP, trace context)
  |    |    |-- requestID.go            # X-Request-ID header generation
  |    |    |-- tracing.go              # New Relic distributed tracing middleware
  |    |
  |    |-- model/
  |    |    |-- model.go                # Base model (CreatedAt, UpdatedAt, DeletedAt for soft delete)
  |    |    |-- user.go                 # User model (email/phone verification, banning, roles)
  |    |    |-- listing.go              # Listing model (JSONB attributes, media URLs, view counter)
  |    |    |-- brand.go                # Brand model (slug, profile/banner, verification)
  |    |    |-- category.go             # Category model (hierarchical, dynamic attributes)
  |    |    |-- session.go              # Session model (hashed token, IP/UA, sliding expiry)
  |    |    |-- message.go              # Message and Conversation models
  |    |    |-- feedback.go             # Feedback model (bug, suggestion, other)
  |    |    |-- reviews.go              # Review model (1-5 rating, images)
  |    |    |-- saved.go                # Saved/bookmarked listing model
  |    |
  |    |-- repository/
  |    |    |-- repository.go           # Repository aggregator struct
  |    |    |-- user.go                 # User queries (insert with brand transaction, get, verify)
  |    |    |-- session.go              # Session queries (create, get user by session, refresh, delete)
  |    |    |-- brand.go                # Brand queries (get by seller, get by ID, update)
  |    |    |-- category.go             # Category queries (CRUD, recursive CTE for attributes)
  |    |    |-- listing.go              # Listing queries (CRUD, dynamic filtering, view increment)
  |    |
  |    |-- router/
  |    |    |-- router.go               # Route definitions and middleware chain assembly
  |    |
  |    |-- server/
  |    |    |-- server.go               # Server struct, initialization, HTTP server, graceful shutdown
  |    |
  |    |-- service/
  |    |    |-- services.go             # Service aggregator struct
  |    |    |-- auth.go                 # Authentication logic (register, login, verify, session management)
  |    |    |-- listing.go              # Listing business logic (CRUD, view tracking, upload signatures)
  |    |    |-- brand.go                # Brand business logic (get, update with image upload)
  |    |    |-- category.go             # Category business logic (CRUD, image management, attribute merging)
  |    |
  |    |-- sqlerr/
  |    |    |-- error.go                # PostgreSQL error code mapping and severity classification
  |    |    |-- handler.go              # pgx/pgconn error to HTTP error conversion
  |    |
  |    |-- validation/
  |         |-- validation.go           # JSON body binding, struct validation, human-readable errors
  |
  |-- pkg/                              # Public types (importable by other packages)
  |    |-- types/
  |         |-- auth.go                 # Auth request/response DTOs
  |         |-- listing.go              # Listing request/response DTOs with filter params
  |         |-- brand.go                # Brand request/response DTOs
  |         |-- category.go             # Category request/response DTOs
  |         |-- file.go                 # File upload metadata type
  |
  |-- static/                           # Static assets
  |    |-- openapi.json                 # OpenAPI 3.0.3 specification (root)
  |    |-- openapi.bundle.json          # Bundled OpenAPI spec
  |    |-- openapi.html                 # Interactive API documentation page
  |    |-- openapi/                     # Split OpenAPI spec files
  |         |-- paths/                  # Path definitions (auth, brands, categories, health, listings)
  |         |-- schemas/                # Schema definitions (brand, listing)
  |
  |-- templates/                        # HTML email templates
  |    |-- emails/
  |         |-- verification.html       # Email verification code template (6-digit code, 10-min expiry)
  |         |-- welcome.html            # Welcome email with feature highlights
  |
  |-- packages/                         # Generated packages (placeholder)
  |    |-- openapi/                     # Generated OpenAPI client (future)
  |    |-- zod/                         # Generated Zod schemas from OpenAPI (future)
  |
  |-- go.mod                            # Go module definition and dependencies
  |-- seed_categories.sql               # Category and attribute seed data
  |-- Taskfile.yml                      # Task runner configuration
```

## Prerequisites

- **Go** 1.25 or later
- **PostgreSQL** 16 -- provided via Docker Compose from the project root
- **Redis** 7 -- provided via Docker Compose from the project root
- **Docker** and **Docker Compose** -- for infrastructure services
- **Task** (optional but recommended) -- task runner (https://taskfile.dev)

### External Service Accounts

| Service    | Required | Purpose                                                           |
| ---------- | -------- | ----------------------------------------------------------------- |
| Cloudinary | Yes      | Image and video upload, storage, transformation, and delivery     |
| Resend     | Yes      | Transactional email delivery (verification codes, welcome emails) |
| New Relic  | No       | Application performance monitoring and distributed tracing        |

## Getting Started

### 1. Start infrastructure services

From the repository root:

```bash
docker-compose up -d
```

This starts PostgreSQL 16 on port `5433` and Redis 7 on port `6379`.

### 2. Set environment variables

Export all required `CAMPUS_CART_*` environment variables (see [Environment Variables](#environment-variables) below). At minimum, you need:

- `CAMPUS_CART_PRIMARY_ENV=development`
- `CAMPUS_CART_SERVER_PORT=<port>`
- `CAMPUS_CART_DB_DSN=postgres://campusCart:campusCart@localhost:5433/campusCart`
- `CAMPUS_CART_REDIS_ADDRESS=localhost:6379`
- `CAMPUS_CART_SERVER_CORS_ALLOWED_ORIGINS=http://localhost:3000`
- `CAMPUS_CART_INTEGRATION_RESEND_API_KEY=<your-resend-key>`
- `CAMPUS_CART_CLOUDINARY_CLOUD_NAME=<your-cloud-name>`
- `CAMPUS_CART_CLOUDINARY_API_KEY=<your-api-key>`
- `CAMPUS_CART_CLOUDINARY_API_SECRET=<your-api-secret>`
- `CAMPUS_CART_AUTH_COOKIE_DOMAIN=localhost`

### 3. Run migrations

```bash
task migrations:up
```

Or, if not using Task, migrations run automatically on server startup in non-development environments.

### 4. Seed categories

```bash
psql -h localhost -p 5433 -U campusCart -d campusCart -f seed_categories.sql
```

This inserts 7 parent categories (Electronics, Books and Stationery, Fashion, Hostel and Room Essentials, Sports and Fitness, Services, Vehicles and Transport), their subcategories, and dynamic category-specific attributes.

### 5. Start the server

```bash
task run
```

Or manually:

```bash
go run ./cmd/campusCart
```

The server will:

1. Load configuration from environment variables.
2. Initialize the logger and optionally New Relic.
3. Run database migrations (non-dev environments).
4. Create database connection pool.
5. Initialize Redis and Cloudinary clients.
6. Wire up repositories, services, handlers, middleware, and routes.
7. Start the HTTP server.
8. Listen for SIGINT/SIGTERM for graceful shutdown (30-second timeout).

## Environment Variables

All environment variables use the prefix `CAMPUS_CART_` and are loaded via the koanf library.

### Primary

| Variable                  | Required | Description                                                 |
| ------------------------- | -------- | ----------------------------------------------------------- |
| `CAMPUS_CART_PRIMARY_ENV` | Yes      | Environment name: `development`, `staging`, or `production` |

### Server

| Variable                                  | Required | Description                                    |
| ----------------------------------------- | -------- | ---------------------------------------------- |
| `CAMPUS_CART_SERVER_PORT`                 | Yes      | HTTP server listening port                     |
| `CAMPUS_CART_SERVER_READ_TIMEOUT`         | No       | HTTP read timeout (default per Go http.Server) |
| `CAMPUS_CART_SERVER_WRITE_TIMEOUT`        | No       | HTTP write timeout                             |
| `CAMPUS_CART_SERVER_IDLE_TIMEOUT`         | No       | HTTP idle connection timeout                   |
| `CAMPUS_CART_SERVER_CORS_ALLOWED_ORIGINS` | Yes      | Comma-separated list of allowed CORS origins   |

### Database

| Variable                            | Required | Description                                     |
| ----------------------------------- | -------- | ----------------------------------------------- |
| `CAMPUS_CART_DB_DSN`                | Yes      | Full PostgreSQL connection string               |
| `CAMPUS_CART_DB_MAX_CONNS`          | No       | Maximum number of connections in the pool       |
| `CAMPUS_CART_DB_MIN_CONNS`          | No       | Minimum number of idle connections maintained   |
| `CAMPUS_CART_DB_MAX_CONN_LIFETIME`  | No       | Maximum total lifetime of a connection          |
| `CAMPUS_CART_DB_MAX_CONN_IDLE_TIME` | No       | Maximum idle time before a connection is closed |

### Redis

| Variable                     | Required | Description                                |
| ---------------------------- | -------- | ------------------------------------------ |
| `CAMPUS_CART_REDIS_ADDRESS`  | Yes      | Redis server address in `host:port` format |
| `CAMPUS_CART_REDIS_PASSWORD` | No       | Redis authentication password              |
| `CAMPUS_CART_REDIS_DB`       | No       | Redis database index (default 0)           |

### Integrations

| Variable                                 | Required | Description                                     |
| ---------------------------------------- | -------- | ----------------------------------------------- |
| `CAMPUS_CART_INTEGRATION_RESEND_API_KEY` | Yes      | Resend API key for sending transactional emails |

### Authentication

| Variable                         | Required | Description                                                            |
| -------------------------------- | -------- | ---------------------------------------------------------------------- |
| `CAMPUS_CART_AUTH_COOKIE_DOMAIN` | Yes      | Domain for the session cookie (e.g., `localhost` or `.yourdomain.com`) |

### Cloudinary

| Variable                            | Required | Description           |
| ----------------------------------- | -------- | --------------------- |
| `CAMPUS_CART_CLOUDINARY_CLOUD_NAME` | Yes      | Cloudinary cloud name |
| `CAMPUS_CART_CLOUDINARY_API_KEY`    | Yes      | Cloudinary API key    |
| `CAMPUS_CART_CLOUDINARY_API_SECRET` | Yes      | Cloudinary API secret |

### Observability

| Variable                                                 | Required | Description                                                      |
| -------------------------------------------------------- | -------- | ---------------------------------------------------------------- |
| `CAMPUS_CART_OBSERVABILITY_NEW_RELIC_LICENSE_KEY`        | No       | New Relic license key (enables APM if set)                       |
| `CAMPUS_CART_OBSERVABILITY_NEW_RELIC_APP_NAME`           | No       | Application name in New Relic dashboard                          |
| `CAMPUS_CART_OBSERVABILITY_LOGGING_LEVEL`                | No       | Log level: trace, debug, info, warn, error (default: info)       |
| `CAMPUS_CART_OBSERVABILITY_LOGGING_FORMAT`               | No       | Log format: json or console (default: json)                      |
| `CAMPUS_CART_OBSERVABILITY_LOGGING_SLOW_QUERY_THRESHOLD` | No       | Duration threshold for logging slow SQL queries (default: 100ms) |
| `CAMPUS_CART_OBSERVABILITY_HEALTH_CHECK_INTERVAL`        | No       | Health check polling interval (default: 30s)                     |

## Development Tasks

The project uses [Task](https://taskfile.dev) as a task runner. Available tasks are defined in `Taskfile.yml`:

| Command                | Description                                                 |
| ---------------------- | ----------------------------------------------------------- |
| `task run`             | Start the Go server with `go run ./cmd/campusCart`          |
| `task tidy`            | Run `go mod tidy` to clean up module dependencies           |
| `task migrations:new`  | Create a new SQL migration file in the migrations directory |
| `task migrations:up`   | Apply all pending database migrations                       |
| `task migrations:down` | Roll back the most recent migration                         |

The `migrations:up` and `migrations:down` tasks use the `CAMPUS_CART_DB_DSN` environment variable to connect to the database.

## Architecture

The server follows a layered architecture with clear separation of concerns:

```
HTTP Request
    |
    v
Router (chi)
    |
    v
Middleware Chain
    |
    v
Handler Layer          -- HTTP concerns: request parsing, response serialization
    |
    v
Service Layer          -- Business logic, orchestration, validation rules
    |
    v
Repository Layer       -- Data access, SQL queries, PostgreSQL interaction
    |
    v
PostgreSQL / Redis / Cloudinary / Resend
```

### Dependency Injection

Dependencies are wired together in `cmd/campusCart/main.go` using constructor injection:

1. **Config** is loaded from environment variables.
2. **Server** struct is created, holding config, logger, database pool, Redis client, and job service.
3. **Repositories** are created with the database pool.
4. **Services** are created with repositories, server context, and external clients.
5. **Handlers** are created with services and server context.
6. **Middleware** is created with services and server context.
7. **Router** is created with handlers and middleware.
8. The HTTP server starts with the assembled router.

There is no dependency injection framework. All wiring is explicit.

### Aggregator Pattern

Each architectural layer uses an aggregator struct to bundle related components:

- `handler.Handlers` aggregates `Health`, `OpenAPI`, `Auth`, `Category`, `Listing`, and `Brand` handlers.
- `service.Services` aggregates `Auth`, `Job`, `Category`, `Listing`, and `Brand` services.
- `repository.Repository` aggregates `Session`, `User`, `Brand`, `Category`, and `Listing` repositories.
- `middleware.Middleware` aggregates `Global`, `ContextEnhancer`, `Tracing`, `Auth`, and `Authorization` middleware.

## Request Lifecycle

A typical authenticated request passes through the following stages:

1. **Request ID** (`requestID.go`): Reads the `X-Request-ID` header or generates a UUID. Stores it in the request context and sets it on the response.

2. **Real IP** (chi built-in): Extracts the real client IP from `X-Forwarded-For` or `X-Real-IP` headers.

3. **CORS** (`global.go`): Validates the `Origin` header against the configured allowed origins. Sets appropriate CORS headers for preflight and actual requests.

4. **New Relic Tracing** (`tracing.go`): Creates a New Relic transaction for the request. Adds custom attributes (method, path, request ID, IP).

5. **Context Enrichment** (`context.go`): Enhances the zerolog logger in the request context with the request ID, client IP, HTTP method, URL path, and New Relic trace/span IDs.

6. **Request Logging** (`global.go`): Logs the start and completion of every request with status code, duration, bytes written. Uses appropriate log levels based on HTTP status (info for 2xx/3xx, warn for 4xx, error for 5xx).

7. **Panic Recovery** (`global.go`): Catches panics from downstream handlers. If the panic value is an `HTTPError`, it is serialized as a JSON error response. Other panics produce a 500 Internal Server Error with the panic value logged.

8. **Security Headers** (`global.go`): Sets `Strict-Transport-Security`, `X-Content-Type-Options: nosniff`, `X-XSS-Protection`, and `X-Frame-Options: DENY`.

9. **Rate Limiting** (chi httprate): Enforces 100 requests per minute per client IP.

10. **Authentication** (`auth.go`, applied to protected routes): Reads the `cc_refresh_token` cookie, hashes it with SHA-256, looks up the session in the database, verifies it has not expired, refreshes the session's last activity timestamp (sliding expiration), and stores the user and brand ID in the request context.

11. **Authorization** (`admin.go`, applied to admin routes): Checks that the authenticated user has the `admin` role. Returns 403 Forbidden otherwise.

12. **Handler Execution** (`base.go`): The generic handler framework:
    - Binds and validates the request body/params using `BindAndValidate`.
    - Executes the handler function.
    - Serializes the response (JSON, No Content, or File).
    - Records New Relic attributes and zerolog durations for the validation and handler phases.

## Configuration

Configuration is managed through the `internal/config` package using the koanf library. All configuration is loaded from environment variables with the `CAMPUS_CART_` prefix.

### Configuration Structure

```go
Config
  |-- Primary          // Environment name
  |-- Server           // Port, timeouts, CORS origins
  |-- Database         // DSN, pool settings (max conns, idle, lifetime)
  |-- Redis            // Address, password, DB index
  |-- Integration      // Resend API key
  |-- Auth             // Cookie domain
  |-- Cloudinary       // Cloud name, API key, API secret
  |-- Observability    // New Relic, logging, health checks
```

### Validation

The configuration struct uses `go-playground/validator` tags. If required fields are missing or invalid, the server fails to start with a descriptive validation error.

### Defaults

| Setting               | Default |
| --------------------- | ------- |
| Log level             | `info`  |
| Log format            | `json`  |
| Slow query threshold  | `100ms` |
| Health check interval | `30s`   |

## Database

### PostgreSQL 16

The server uses pgx v5 with the `pgxpool` connection pool. The pool is configured with:

- Maximum and minimum connection counts
- Connection lifetime and idle time limits
- DSN constructed from configuration

### Tracing

In development and staging environments, the database connection includes multi-tracer support:

- **New Relic tracer**: Records SQL queries as datastore segments in New Relic.
- **zerolog tracer**: Logs SQL queries with the configured slow query threshold.

### Connection Validation

On startup, the server validates the database connection by executing a test query. If the connection fails, the server aborts.

### Extensions

The database uses two PostgreSQL extensions, created in the first migration:

- `uuid-ossp`: Provides `uuid_generate_v4()` for primary key generation.
- `citext`: Provides case-insensitive text type for email fields.

### Custom Enums

| Enum                | Values                 |
| ------------------- | ---------------------- |
| `user_role`         | admin, user            |
| `listing_condition` | new, used, refurbished |
| `media_type`        | image, video           |
| `feedback_type`     | suggestion, bug, other |

### Trigger Function

A `set_updated_at()` trigger function is created in the first migration and applied to all mutable tables. It automatically sets the `updated_at` column to `NOW()` on every row update.

## Migrations

Migrations are managed using tern v2. Migration SQL files are embedded into the Go binary using Go's `embed.FS` at compile time, ensuring they are always available regardless of the deployment environment.

### Migration Files

| File                            | Description                                                                |
| ------------------------------- | -------------------------------------------------------------------------- |
| `001_setup.sql`                 | Extensions (uuid-ossp, citext), enums, `set_updated_at()` trigger function |
| `002_create-user-table.sql`     | Users table with email/phone verification, banning, soft delete            |
| `003_create-session-table.sql`  | Sessions with hashed refresh tokens, IP/UA tracking, expiry                |
| `004_create-category-table.sql` | Hierarchical categories with dynamic attributes (JSONB options)            |
| `005_create-brand-table.sql`    | Seller brands with slug, profile/banner URLs, verification                 |
| `006_create-listing-table.sql`  | Listings with JSONB attributes, media URL arrays, view counter             |
| `007_create-saved-table.sql`    | Saved/bookmarked listings (unique per user+listing)                        |
| `008_create-review-table.sql`   | Reviews with 1-5 rating (unique per listing+reviewer)                      |
| `009_create-message-table.sql`  | Conversations and messages with media support                              |
| `010_create-feedback-table.sql` | User feedback (suggestion, bug, other) with admin notes                    |

### Running Migrations

```bash
# Apply all pending migrations
task migrations:up

# Roll back the last migration
task migrations:down

# Create a new migration file
task migrations:new
```

In non-development environments, migrations are applied automatically on server startup before the HTTP server begins accepting requests.

## Seeding

The `seed_categories.sql` file populates the database with initial category data. It uses an upsert pattern (`INSERT ... ON CONFLICT DO UPDATE`) so it can be run multiple times safely.

### Parent Categories

1. Electronics
2. Books and Stationery
3. Fashion
4. Hostel and Room Essentials
5. Sports and Fitness
6. Services
7. Vehicles and Transport

### Subcategories (Examples)

- Electronics: Phones and Tablets, Laptops and Computers, Audio and Headphones, Chargers and Accessories, Gaming
- Fashion: Clothing, Shoes, Bags and Backpacks, Accessories and Jewelry
- Services: Tutoring, Printing and Design, Photography, Tech Repair, Delivery and Errands

### Dynamic Attributes

Each subcategory has typed attributes that listings in that category must provide. For example:

- Phones and Tablets: condition (enum), brand (text), model (text), storage_capacity (text), battery_health (text)
- Laptops: condition (enum), brand (text), model (text), ram (text), storage (text), processor (text)
- Clothing: condition (enum), size (text), gender (enum: male/female/unisex), material (text)
- Tutoring: subject (text), level (text), mode (enum: online/in-person/both), rate_type (enum: hourly/fixed/negotiable)

Attribute types include `text`, `number`, `boolean`, and `enum` (with predefined options stored as JSONB).

## Authentication and Sessions

### Registration Flow

1. User submits first name, last name, username, email, and password.
2. Email is validated to be a student email (`@st.ug.edu.gh` suffix required).
3. Email uniqueness is checked.
4. Password is hashed with bcrypt.
5. User is created in a database transaction that also creates a default brand (with a slug derived from the username).
6. A 6-digit verification code is generated with a 10-minute expiry.
7. A verification email is dispatched via the async job queue (Asynq).

### Email Verification

1. User submits email and 6-digit verification code.
2. Code is validated against the stored code and expiry.
3. User's `email_verified` flag is set to true.
4. A session is created and the `cc_refresh_token` cookie is set.
5. A welcome email is dispatched via the async job queue.

### Login Flow

1. User submits email and password.
2. User is looked up by email.
3. Password is compared against the bcrypt hash.
4. A new session is created with a 32-byte random token (base64-encoded).
5. The token is hashed with SHA-256 and stored in the database along with the user agent and IP address.
6. The raw token is set as the `cc_refresh_token` cookie.

### Session Cookie

| Property | Value                                           |
| -------- | ----------------------------------------------- |
| Name     | `cc_refresh_token`                              |
| Max Age  | 7 days                                          |
| HttpOnly | true                                            |
| SameSite | Lax                                             |
| Secure   | true in production, false in development        |
| Domain   | Configured via `CAMPUS_CART_AUTH_COOKIE_DOMAIN` |
| Path     | `/`                                             |

### Sliding Expiration

Every authenticated request refreshes the session by updating its `last_activity_at` timestamp and extending the expiry by 7 days from the current time. This means active users never need to re-authenticate as long as they make at least one request within every 7-day window.

### Logout

The session is deleted from the database and the cookie is cleared (max age set to -1).

## API Endpoints

### Public Endpoints

| Method | Path                          | Handler                  | Description                                        |
| ------ | ----------------------------- | ------------------------ | -------------------------------------------------- |
| `POST` | `/auth/register`              | `Auth.Register`          | Create a new user account                          |
| `POST` | `/auth/login`                 | `Auth.Login`             | Authenticate with email and password               |
| `POST` | `/auth/verify-email`          | `Auth.VerifyEmail`       | Verify email with 6-digit code                     |
| `GET`  | `/categories`                 | `Category.GetAll`        | List all categories                                |
| `GET`  | `/categories/{id}`            | `Category.GetByID`       | Get a single category by ID                        |
| `GET`  | `/categories/{id}/attributes` | `Category.GetAttributes` | Get attributes for a category (merged with parent) |
| `GET`  | `/listings`                   | `Listing.List`           | Browse listings with filters and pagination        |
| `GET`  | `/listings/{id}`              | `Listing.Get`            | Get a single listing (increments view count)       |
| `GET`  | `/docs`                       | `OpenAPI.Docs`           | Interactive OpenAPI documentation page             |

### Authenticated Endpoints

| Method   | Path                         | Handler                   | Description                                      |
| -------- | ---------------------------- | ------------------------- | ------------------------------------------------ |
| `POST`   | `/auth/logout`               | `Auth.Logout`             | End the current session                          |
| `GET`    | `/auth/me`                   | `Auth.GetCurrentUser`     | Get the authenticated user's profile             |
| `GET`    | `/brands/me`                 | `Brand.GetCurrent`        | Get the authenticated user's brand               |
| `PATCH`  | `/brands/me`                 | `Brand.Update`            | Update brand (JSON or multipart with images)     |
| `POST`   | `/listings`                  | `Listing.Create`          | Create a new listing                             |
| `PATCH`  | `/listings/{id}`             | `Listing.Update`          | Update an existing listing                       |
| `DELETE` | `/listings/{id}`             | `Listing.Delete`          | Soft-delete a listing (brand ownership required) |
| `POST`   | `/listings/upload-signature` | `Listing.UploadSignature` | Get a signed Cloudinary upload URL               |

### Admin Endpoints

| Method   | Path               | Handler           | Description                                  |
| -------- | ------------------ | ----------------- | -------------------------------------------- |
| `POST`   | `/categories`      | `Category.Create` | Create a new category (multipart with image) |
| `PATCH`  | `/categories/{id}` | `Category.Update` | Update a category (multipart with image)     |
| `DELETE` | `/categories/{id}` | `Category.Delete` | Soft-delete a category                       |

### Listing Filters

The `GET /listings` endpoint supports the following query parameters:

| Parameter     | Type    | Description                                                           |
| ------------- | ------- | --------------------------------------------------------------------- |
| `category_id` | UUID    | Filter by category (includes descendant categories via recursive CTE) |
| `brand_id`    | UUID    | Filter by brand                                                       |
| `search`      | string  | Full-text search across listing title and description (ILIKE)         |
| `min_price`   | number  | Minimum price filter                                                  |
| `max_price`   | number  | Maximum price filter                                                  |
| `condition`   | string  | Filter by condition (new, used, refurbished)                          |
| `limit`       | integer | Number of results per page (default and maximum enforced by service)  |
| `offset`      | integer | Pagination offset                                                     |

Results are ordered with promoted listings first, then by creation date descending.

### Health Check

The `GET /health` endpoint checks connectivity to:

- PostgreSQL (via `SELECT 1`)
- Redis (via `PING`)

Each check reports its status and response time. If any check fails, a New Relic custom event is recorded.

## Middleware Chain

The middleware chain is assembled in `internal/router/router.go` and applied in the following order:

| Order | Middleware        | File           | Description                                  |
| ----- | ----------------- | -------------- | -------------------------------------------- |
| 1     | Request ID        | `requestID.go` | Generate or propagate `X-Request-ID`         |
| 2     | Real IP           | chi built-in   | Extract client IP from proxy headers         |
| 3     | CORS              | `global.go`    | Cross-origin request handling                |
| 4     | New Relic         | `tracing.go`   | Distributed tracing and transaction creation |
| 5     | Context Enricher  | `context.go`   | Add structured fields to the request logger  |
| 6     | Request Logger    | `global.go`    | Log request start/completion with timing     |
| 7     | Recoverer         | `global.go`    | Panic recovery with error response           |
| 8     | Security Headers  | `global.go`    | HSTS, nosniff, XSS protection, frame denial  |
| 9     | Rate Limiter      | chi httprate   | 100 requests/minute per client               |
| 10    | Auth (per-route)  | `auth.go`      | Session validation and sliding expiration    |
| 11    | Admin (per-route) | `admin.go`     | Admin role check                             |

## Error Handling

The server has a comprehensive error handling system with three tiers:

### 1. HTTP Errors (`internal/err/`)

The `HTTPError` struct represents all application-level errors:

```go
type HTTPError struct {
    Code     string              // Machine-readable error code
    Message  string              // Human-readable description
    Status   int                 // HTTP status code
    Override bool                // Whether to use custom status code
    Errors   map[string]string   // Field-level validation errors
    Action   *Action             // Optional redirect action
}
```

Factory functions: `NewUnauthorizedError`, `NewForbiddenError`, `NewBadRequestError`, `NewNotFoundError`, `NewInternalServerError`, `ValidationError`.

### 2. SQL Error Translation (`internal/sqlerr/`)

PostgreSQL errors from pgx are automatically translated to user-friendly HTTP errors:

| PostgreSQL Code               | HTTP Status     | Example                       |
| ----------------------------- | --------------- | ----------------------------- |
| 23505 (unique_violation)      | 409 Conflict    | "Email already exists"        |
| 23503 (foreign_key_violation) | 400 Bad Request | "Referenced entity not found" |
| 23502 (not_null_violation)    | 400 Bad Request | "Required field missing"      |
| 23514 (check_violation)       | 400 Bad Request | "Constraint check failed"     |
| 40P01 (deadlock_detected)     | 409 Conflict    | "Transaction conflict, retry" |

The error handler auto-generates meaningful error codes and entity names from the PostgreSQL constraint and table names.

### 3. Panic Recovery (`internal/middleware/global.go`)

If a handler panics with an `HTTPError`, it is caught and serialized as a normal error response. Other panics produce a 500 Internal Server Error with the panic details logged.

## Validation

Request validation is handled by the `internal/validation` package.

### Bind and Validate

The `BindAndValidate` function:

1. Decodes the JSON request body into the target struct.
2. Calls the `Validate()` method if the struct implements the `Validatable` interface (for custom validation rules).
3. Runs `go-playground/validator` struct tag validation.
4. Translates validation errors into human-readable messages.

### Validation Tags

Supported tags and their human-readable messages:

| Tag        | Message                                                            |
| ---------- | ------------------------------------------------------------------ |
| `required` | "is required"                                                      |
| `min`      | "must be at least {param}" / "must be at least {param} characters" |
| `max`      | "must be at most {param}" / "must be at most {param} characters"   |
| `oneof`    | "must be one of: {param}"                                          |
| `email`    | "must be a valid email address"                                    |
| `uuid`     | "must be a valid UUID"                                             |
| `url`      | "must be a valid URL"                                              |
| `gte`      | "must be greater than or equal to {param}"                         |
| `lte`      | "must be less than or equal to {param}"                            |
| `alphanum` | "must contain only alphanumeric characters"                        |

Field names are automatically converted from Go struct field names to JSON-style snake_case.

## File Uploads

File uploads are handled through Cloudinary using two approaches:

### 1. Server-Side Upload

Used for category images and brand profile/banner images. The handler:

1. Parses the multipart form data.
2. Validates file size (2MB for categories, 5MB for brands) and MIME type (images only).
3. Uploads the file to Cloudinary via the `lib/file` client.
4. Stores the returned URL in the database.

### 2. Direct-to-Cloudinary Upload (Client-Side)

Used for listing images. The flow:

1. Client requests a signed upload URL from `POST /listings/upload-signature`.
2. Server generates a SHA-1 signed parameter set (timestamp, tags, context, folder, API key, upload preset).
3. Client uploads directly to Cloudinary using the signed parameters.
4. Client includes the returned Cloudinary URLs in the listing creation/update request.

This approach avoids routing large media files through the application server.

### Cloudinary Operations

The `lib/file/client.go` provides:

- `UploadFile(file, tags, metadata)` -- server-side upload
- `DeleteFile(publicID)` -- delete by public ID
- `GenerateDirectUpload(file)` -- generate signed upload parameters for client-side upload

## Email System

Transactional emails are sent through the Resend API using HTML templates.

### Email Types

| Email             | Template            | Trigger                       |
| ----------------- | ------------------- | ----------------------------- |
| Verification Code | `verification.html` | User registration             |
| Welcome           | `welcome.html`      | Successful email verification |

### Template Rendering

Templates are stored in `templates/emails/` as HTML files with Go `html/template` placeholders. The email client renders these templates with the provided data and sends the resulting HTML via Resend.

### Verification Email

The verification email contains a 6-digit code displayed prominently in a branded card layout (using the `#4B3EC4` brand color). The code expires after 10 minutes.

### Welcome Email

The welcome email includes a feature overview and a call-to-action button directing the user to the platform.

## Background Jobs

Background job processing is handled by Asynq, a distributed task queue built on Redis.

### Queue Configuration

| Queue      | Priority | Use Case                  |
| ---------- | -------- | ------------------------- |
| `critical` | Highest  | Time-sensitive operations |
| `default`  | Medium   | Standard background tasks |
| `low`      | Lowest   | Non-urgent operations     |

The worker pool processes up to 10 concurrent tasks.

### Defined Job Types

| Type                 | Description                                     |
| -------------------- | ----------------------------------------------- |
| `email:welcome`      | Send welcome email after verification           |
| `email:verification` | Send verification code email after registration |

### Job Lifecycle

1. A service creates a task payload (e.g., email recipient and template data).
2. The task is enqueued to the appropriate queue via the job service.
3. The Asynq worker picks up the task and executes the registered handler.
4. The handler (e.g., `HandleWelcomeEmail`) deserializes the payload and calls the email client.

## Logging

Logging is implemented using zerolog, a zero-allocation JSON logger.

### Output Formats

| Environment | Format  | Features                                                      |
| ----------- | ------- | ------------------------------------------------------------- |
| Development | Console | Human-readable, colorized, timestamp formatting, stack traces |
| Production  | JSON    | Structured, machine-parseable, New Relic log forwarding       |

### Log Levels

The log level is configurable via `CAMPUS_CART_OBSERVABILITY_LOGGING_LEVEL`. Available levels: trace, debug, info, warn, error, fatal, panic.

### Request Logging

Every HTTP request is logged with:

- Request ID
- Client IP address
- HTTP method and path
- Response status code
- Response duration
- Bytes written

Log level is determined by response status: info for 2xx/3xx, warn for 4xx, error for 5xx.

### SQL Query Logging

SQL queries executed through pgx are logged with:

- Query text (truncated for readability)
- Execution duration
- Slow query warnings when duration exceeds the configured threshold (default 100ms)

## Observability

### New Relic APM

When a New Relic license key is configured, the server provides:

- **Transaction tracing**: Every HTTP request is a New Relic transaction with custom attributes (method, path, request ID, user ID, role).
- **Datastore segments**: PostgreSQL queries are recorded as datastore operations.
- **Custom events**: Health check failures trigger custom New Relic events.
- **Log forwarding**: In production, zerolog entries are forwarded to New Relic Logs.

### Health Checks

The `/health` endpoint reports:

- PostgreSQL connectivity and response time.
- Redis connectivity and response time.
- Overall system status.

The health check interval is configurable (default 30 seconds).

## Domain Models

All models include a base `Model` struct with `CreatedAt`, `UpdatedAt`, and `DeletedAt` (nullable, for soft delete) timestamps.

### User

Fields: ID, first name, last name, username, email (citext), phone, password hash, role (admin/user), email verified flag, verification code and expiry, phone verification code and expiry, banned flag, active flag.

### Brand

Fields: ID, seller ID (FK to users), name, slug (unique), description, profile URL, banner URL, is verified flag.

A default brand is created automatically for every new user during registration.

### Listing

Fields: ID, brand ID (FK to brands), category ID (FK to categories), title, description, price (numeric), condition (new/used/refurbished), attributes (JSONB), image URLs (text array, non-empty constraint), video URLs (text array), views count, promoted flag.

### Category

Fields: ID, name, slug (unique), description, image URL, parent ID (self-referential FK for hierarchy).

### Category Attribute

Fields: ID, category ID (FK), name, attribute type (text/number/boolean/enum), is required flag, options (JSONB for enum values).

### Session

Fields: ID, user ID (FK), refresh token hash, user agent, IP address, expires at, last activity at.

### Conversation

Fields: ID, buyer ID, seller ID, listing ID, last message at.

### Message

Fields: ID, conversation ID, sender ID, content, media URL, media type (image/video), is read flag, read at.

### Feedback

Fields: ID, user ID, type (suggestion/bug/other), subject, message, status, admin notes.

### Review

Fields: ID, listing ID, reviewer ID, rating (1-5), comment, image URLs (text array).

### Saved

Fields: ID, user ID, listing ID. Unique constraint on (user_id, listing_id).

## Repository Layer

The repository layer handles all direct database interactions using pgx v5. Each repository struct takes a `*pgxpool.Pool` in its constructor.

### Key Patterns

- **Soft Delete**: Queries filter on `deleted_at IS NULL` by default.
- **Transactions**: Multi-table operations (e.g., user + brand creation during registration) use `pgx.BeginTxFunc`.
- **COALESCE Updates**: The `UpdateBrand` and similar functions use `COALESCE(NULLIF($param, ''), existing_column)` to only update fields that are provided (non-empty), preserving existing values for omitted fields.
- **Recursive CTE**: Category attribute queries use recursive Common Table Expressions to traverse the category hierarchy and merge parent attributes with child attributes (child attributes win on name collision via `DISTINCT ON`).
- **Dynamic Query Building**: The listing `List()` function constructs SQL queries dynamically based on provided filter parameters, using parameterized queries to prevent SQL injection.

## Service Layer

The service layer contains all business logic. Each service struct receives its dependencies (repositories, server context, external clients) via constructor injection.

### Key Service Behaviors

**Auth Service**

- Validates student email domain (`@st.ug.edu.gh`).
- Generates cryptographically random session tokens (32 bytes, base64-encoded).
- Generates 6-digit numeric verification codes.
- Orchestrates user creation, session management, and async email dispatch.

**Listing Service**

- Enforces default and maximum pagination limits.
- Separates "get" (read-only) from "view" (increments view counter) operations.
- Validates brand ownership before allowing updates or deletes.
- Generates signed Cloudinary upload parameters for client-side uploads.

**Brand Service**

- Handles optional image uploads to Cloudinary during brand updates.
- Supports both JSON and multipart form-data request formats.

**Category Service**

- Manages image uploads and deletions during category lifecycle.
- Generates URL-friendly slugs from category names.
- Merges category attributes with parent attributes using deduplication (child attributes take precedence).

## Handler Layer

Handlers are the HTTP interface layer. They use a generic handler framework defined in `internal/handler/base.go`.

### Generic Handler Framework

Three handler patterns are provided:

1. **`Handle[Req, Res]`**: For endpoints that return JSON responses. Binds and validates the request, executes the handler function, and serializes the response as JSON.

2. **`HandleNoContent[Req]`**: For endpoints that return 204 No Content (e.g., delete operations). Binds and validates the request, executes the handler, and sends an empty response.

3. **`HandleFile[Req]`**: For endpoints that return file content (e.g., serving static files or generated content).

Each pattern includes:

- Automatic request validation via `BindAndValidate`.
- New Relic attribute recording (validation duration, handler duration).
- Structured logging of operation timings.
- Consistent error handling and response serialization.

### Multipart Form Handling

Brand update and category create/update endpoints support multipart form data for file uploads. The handlers:

1. Parse the multipart form with a size limit.
2. Validate file MIME types (image-only).
3. Pass file headers to the service layer for upload.

## OpenAPI Documentation

The API is documented using the OpenAPI 3.0.3 specification. The spec files are located in `static/openapi/`:

- `openapi.json`: Root spec file with server info and component references.
- `openapi.bundle.json`: Fully bundled spec for standalone use.
- `openapi/paths/`: Individual path definition files (auth, brands, categories, health, listings).
- `openapi/schemas/`: Reusable schema definitions (brand, listing).

An interactive documentation page is served at `GET /docs` using the bundled HTML file at `static/openapi.html`.

### Security Scheme

The API uses cookie-based authentication. The OpenAPI spec defines a `sessionCookie` security scheme of type `apiKey` in the cookie location, referencing the `cc_refresh_token` cookie.

## Uploader CLI Tool

The `cmd/uploader/main.go` provides a standalone CLI tool for uploading files to Cloudinary through the server's signed upload endpoint.

### Usage

```bash
go run ./cmd/uploader/main.go
```

The tool:

1. Requests a signed upload URL from the server's `/listings/upload-signature` endpoint.
2. Constructs a multipart upload request with the file and signed parameters.
3. Uploads directly to Cloudinary.
4. Prints the response (including the uploaded file's URL).

This is useful for seeding media assets during development or testing the upload pipeline outside the browser.

## Code Conventions

### Project Layout

The server follows standard Go project layout conventions:

- `cmd/`: Application entry points. Each subdirectory is a separate binary.
- `internal/`: Private packages not importable by external code.
- `pkg/`: Public packages (request/response types) that could be imported by clients.

### Error Handling

- Never return raw errors to clients. Always wrap errors using the `err` package factories.
- Database errors from pgx should pass through the `sqlerr` handler for automatic translation.
- Use `panic(httpErr)` for early returns in handlers; the recovery middleware converts these to HTTP responses.

### Naming

- Repository methods use SQL-style prefixes: `Insert`, `Select`, `Update`, `Delete`, `Get`.
- Service methods use domain-style names: `Login`, `Register`, `CreateListing`, `ViewListing`.
- Handler methods match the HTTP operation: `Create`, `Get`, `List`, `Update`, `Delete`.

### Database Queries

- Use parameterized queries exclusively. Never interpolate values into SQL strings.
- Use `COALESCE(NULLIF(...), ...)` for partial updates.
- Use soft delete (`deleted_at IS NOT NULL`) rather than physical deletion.
- Include `deleted_at IS NULL` in all SELECT queries by default.

### Context Usage

- User information is stored in the request context by the auth middleware.
- The logger is stored in the request context by the context enrichment middleware.
- Brand ID is stored alongside user info for ownership checks.

## Deployment

### Build

```bash
cd apps/server
go build -o campuscart ./cmd/campusCart
```

The resulting binary is fully self-contained. Migrations and email templates are embedded at compile time.

### Runtime Requirements

- PostgreSQL 16 database (accessible via DSN)
- Redis 7 instance (for sessions, caching, and job queues)
- All required `CAMPUS_CART_*` environment variables set

### Graceful Shutdown

The server listens for `SIGINT` and `SIGTERM` signals. On receiving either:

1. The HTTP server stops accepting new connections (30-second timeout for in-flight requests).
2. The database connection pool is closed.
3. The Asynq job service is shut down.
4. The Redis connection is closed.

### Production Considerations

- Set `CAMPUS_CART_PRIMARY_ENV=production` to enable:
  - Secure session cookies
  - JSON log output
  - New Relic log forwarding
  - Automatic migration execution on startup
- Configure `CAMPUS_CART_SERVER_CORS_ALLOWED_ORIGINS` to include only your production client domain.
- Set `CAMPUS_CART_AUTH_COOKIE_DOMAIN` to your production domain (e.g., `.campuscart.com`).
- Tune database pool settings based on expected traffic and PostgreSQL `max_connections`.
- Monitor health via the `/health` endpoint and New Relic dashboards.
