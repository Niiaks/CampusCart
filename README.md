# CampusCart

CampusCart is a web-based student marketplace platform designed to connect student sellers with student buyers within campus communities. It provides a centralized platform where student sellers can list their products and services, and student buyers can browse, discover, and communicate with sellers -- all without leaving campus.

**Author:** Nii Akwei Pappoe
**Started:** February 2026

---

## Table of Contents

- [Problem Statement](#problem-statement)
- [Vision](#vision)
- [Target Use Cases](#target-use-cases)
- [Architecture Overview](#architecture-overview)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Infrastructure Services](#infrastructure-services)
- [Environment Variables](#environment-variables)
- [Development Workflow](#development-workflow)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [MVP Features](#mvp-features)
- [Post-MVP Roadmap](#post-mvp-roadmap)
- [Success Metrics](#success-metrics)
- [Contributing](#contributing)
- [License](#license)

---

## Problem Statement

Student sellers currently rely on social media platforms such as WhatsApp and Telegram groups to promote and sell their goods, products, or services to potential customers. Buyers must navigate through these fragmented channels to find what they need on campus. This process is cumbersome, inefficient, and lacks any centralized discovery mechanism.

Existing alternative marketplace solutions (such as Jiji and Tonaton) are designed for the general public, not tailored for student communities. They suffer from high fraud rates and do not offer the kind of focused, visually appealing experience that would retain student users.

## Vision

CampusCart aims to become Ghana's largest student marketplace network by offering a web application that enables student sellers to list their products or services in one centralized place. This increases visibility for sellers and makes it faster and easier for student buyers to find what they need on campus.

There is currently no student-only marketplace, which creates a significant opportunity for CampusCart to stand out and compete with general-purpose marketplace competitors. Similar apps have demonstrated viable monetization paths, and CampusCart intends to follow suit once a critical mass of 500+ active sellers is reached.

## Target Use Cases

1. **Sellers** want to list their goods and services for fellow students and others to see and purchase.
2. **Buyers** want to browse goods and services offered by fellow students so they can make purchases without going outside campus.

## Architecture Overview

CampusCart is organized as a monorepo containing two primary applications:

```
campusCart/
  docker-compose.yml        # Infrastructure services (PostgreSQL, Redis)
  apps/
    client/                 # Next.js frontend application
    server/                 # Go backend API server
```

The system follows a client-server architecture:

- **Client**: A Next.js 16 application using React 19, Tailwind CSS 4, and shadcn/ui components. It communicates with the backend API via Axios with cookie-based session authentication. React Query manages server state, caching, and synchronization.
- **Server**: A Go API server built with the chi router framework. It uses PostgreSQL 16 as the primary datastore, Redis 7 for caching and background job queues (via Asynq), Cloudinary for media storage, and Resend for transactional emails. Observability is handled through New Relic APM and zerolog structured logging.
- **Infrastructure**: PostgreSQL and Redis run as Docker containers defined in the root `docker-compose.yml`.

### Request Flow

1. The client sends HTTP requests to the Go API server.
2. Requests pass through a middleware chain: Request ID generation, IP extraction, CORS handling, New Relic tracing, context enrichment, structured request logging, panic recovery, and security headers.
3. Rate limiting is enforced at 100 requests per minute per client.
4. Authentication uses cookie-based sessions with a 7-day sliding expiration window. The session token is stored as an HttpOnly cookie (`cc_refresh_token`).
5. The server processes requests through handler, service, and repository layers.
6. PostgreSQL handles persistence, Redis handles session caching and job queues, Cloudinary handles file uploads, and Resend handles email delivery.

### Data Flow for Key Operations

- **Registration**: Email uniqueness check, student email validation (`@st.ug.edu.gh`), password hashing (bcrypt), user + default brand creation (in a transaction), verification email dispatched via async job queue.
- **Listing Creation**: Authenticated seller creates a listing with images uploaded directly to Cloudinary (client gets a signed upload URL from the server), listing metadata stored in PostgreSQL with JSONB attributes.
- **Browsing**: Public endpoint supports filtering by category (with hierarchical descendant resolution), brand, search term (ILIKE), price range, condition, and pagination. Promoted listings appear first.

## Tech Stack

### Frontend

| Technology             | Purpose                                            |
| ---------------------- | -------------------------------------------------- |
| Next.js 16             | React framework with App Router, server components |
| React 19               | UI library                                         |
| TypeScript 5           | Type safety                                        |
| Tailwind CSS 4         | Utility-first CSS framework                        |
| shadcn/ui (New York)   | Component library built on Radix UI primitives     |
| TanStack React Query 5 | Server state management, caching, synchronization  |
| Axios                  | HTTP client with interceptors                      |
| Zod 4                  | Schema validation for forms and API payloads       |
| React Hook Form 7      | Form state management                              |
| Sonner                 | Toast notifications                                |
| Lucide React           | Icon library                                       |

### Backend

| Technology                  | Purpose                                             |
| --------------------------- | --------------------------------------------------- |
| Go 1.25                     | Server programming language                         |
| chi v5                      | Lightweight HTTP router with middleware support     |
| pgx v5                      | PostgreSQL driver with connection pooling           |
| tern v2                     | SQL migration management (embedded migrations)      |
| Redis v9 (go-redis)         | Caching, session store, job queue backend           |
| Asynq                       | Distributed task queue built on Redis               |
| Cloudinary v2               | Image and video upload, storage, and transformation |
| Resend v2                   | Transactional email delivery                        |
| zerolog                     | Structured JSON logging                             |
| New Relic Go Agent          | Application performance monitoring and tracing      |
| go-playground/validator v10 | Struct validation with custom error messages        |
| bcrypt                      | Password hashing                                    |
| koanf                       | Configuration management from environment variables |

### Infrastructure

| Technology              | Purpose                                       |
| ----------------------- | --------------------------------------------- |
| PostgreSQL 16 (Alpine)  | Primary relational database                   |
| Redis 7 (Alpine)        | Caching, job queues, session storage          |
| Docker & Docker Compose | Container orchestration for local development |

## Project Structure

```
campusCart/
 |-- docker-compose.yml              # PostgreSQL 16 + Redis 7 containers
 |-- README.md                       # This file
 |-- apps/
      |-- client/                    # Next.js frontend (see apps/client/README.md)
      |    |-- app/                  # App Router pages and layouts
      |    |-- common/               # Shared types, schemas, utilities
      |    |-- components/           # Reusable UI components (shadcn/ui)
      |    |-- constants/            # Application constants
      |    |-- hooks/                # Custom React hooks
      |    |-- lib/                  # Library utilities
      |    |-- public/               # Static assets
      |    |-- services/             # API service layer (Axios calls)
      |
      |-- server/                    # Go backend API (see apps/server/README.md)
           |-- cmd/                  # Application entry points
           |-- internal/             # Private application code
           |    |-- config/          # Configuration loading and validation
           |    |-- database/        # Database connection, migrations, pooling
           |    |-- err/             # HTTP error types and factories
           |    |-- handler/         # HTTP handlers (controllers)
           |    |-- lib/             # Internal libraries (email, file, jobs)
           |    |-- logger/          # Structured logging setup
           |    |-- middleware/      # HTTP middleware chain
           |    |-- model/           # Domain models
           |    |-- repository/      # Data access layer
           |    |-- router/          # Route definitions
           |    |-- server/          # Server initialization and lifecycle
           |    |-- service/         # Business logic layer
           |    |-- sqlerr/          # PostgreSQL error translation
           |    |-- validation/      # Request validation
           |-- pkg/                  # Public types (request/response DTOs)
           |-- static/              # OpenAPI spec and docs
           |-- templates/           # Email HTML templates
```

## Prerequisites

Before running CampusCart locally, ensure you have the following installed:

- **Docker** and **Docker Compose** -- for running PostgreSQL and Redis
- **Node.js** (v20 or later) and **npm** -- for the client application
- **Go** (1.25 or later) -- for the server application
- **Task** (Taskfile runner) -- for running server development tasks (optional but recommended; install from https://taskfile.dev)

### External Service Accounts

The following third-party services require accounts and API keys:

- **Cloudinary** -- for image and video upload and storage
- **Resend** -- for transactional email delivery (verification codes, welcome emails)
- **New Relic** (optional) -- for application performance monitoring

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/Niiaks/campusCart.git
cd campusCart
```

### 2. Start Infrastructure Services

```bash
docker-compose up -d
```

This starts:

- **PostgreSQL 16** on port `5433` (mapped to container port `5432`) with database `campusCart`, user `campusCart`, and scram-sha-256 authentication.
- **Redis 7** on port `6379`.

Both services include health checks and persistent data volumes (`postgres_data`, `redis_data`).

### 3. Set Up the Server

```bash
cd apps/server
```

Copy the required environment variables (see [Environment Variables](#environment-variables) below) and configure them for your local setup.

Run database migrations and start the server:

```bash
# Using Taskfile (recommended)
task run

# Or manually
go run ./cmd/campusCart
```

The server runs migrations automatically on startup (in non-development environments). In development, use:

```bash
task migrations:up
```

### 4. Seed the Database

After migrations are applied, seed the category data:

```bash
psql -h localhost -p 5433 -U campusCart -d campusCart -f seed_categories.sql
```

This populates 7 parent categories (Electronics, Books and Stationery, Fashion, Hostel and Room Essentials, Sports and Fitness, Services, Vehicles and Transport) along with their subcategories and dynamic attributes.

### 5. Set Up the Client

```bash
cd apps/client
npm install
```

Create a `.env.local` file with:

```
NEXT_PUBLIC_API_URL=http://localhost:<server-port>
```

Start the development server:

```bash
npm run dev
```

The client application will be available at `http://localhost:3000`.

## Infrastructure Services

### Docker Compose Configuration

The `docker-compose.yml` at the project root defines two services:

**PostgreSQL 16 (Alpine)**

- Container name: `campusCart-postgres`
- Port mapping: `5433:5432`
- Database: `campusCart`
- User: `campusCart`
- Authentication: `scram-sha-256`
- Initialization: Uses `POSTGRES_HOST_AUTH_METHOD` environment variable
- Health check: `pg_isready` every 5 seconds
- Volume: `postgres_data` (persistent)

**Redis 7 (Alpine)**

- Container name: `campusCart-redis`
- Port mapping: `6379:6379`
- Health check: `redis-cli ping` every 5 seconds
- Volume: `redis_data` (persistent)

## Environment Variables

### Server

All server environment variables use the prefix `CAMPUS_CART_` and are loaded via the koanf configuration library.

| Variable                                          | Description                                         |
| ------------------------------------------------- | --------------------------------------------------- |
| `CAMPUS_CART_PRIMARY_ENV`                         | Environment name (development, staging, production) |
| `CAMPUS_CART_SERVER_PORT`                         | HTTP server port                                    |
| `CAMPUS_CART_SERVER_READ_TIMEOUT`                 | HTTP read timeout                                   |
| `CAMPUS_CART_SERVER_WRITE_TIMEOUT`                | HTTP write timeout                                  |
| `CAMPUS_CART_SERVER_IDLE_TIMEOUT`                 | HTTP idle timeout                                   |
| `CAMPUS_CART_SERVER_CORS_ALLOWED_ORIGINS`         | Comma-separated list of allowed CORS origins        |
| `CAMPUS_CART_DB_DSN`                              | PostgreSQL connection string                        |
| `CAMPUS_CART_DB_MAX_CONNS`                        | Maximum database connections in pool                |
| `CAMPUS_CART_DB_MIN_CONNS`                        | Minimum idle database connections                   |
| `CAMPUS_CART_DB_MAX_CONN_LIFETIME`                | Maximum connection lifetime                         |
| `CAMPUS_CART_DB_MAX_CONN_IDLE_TIME`               | Maximum connection idle time                        |
| `CAMPUS_CART_REDIS_ADDRESS`                       | Redis server address (host:port)                    |
| `CAMPUS_CART_REDIS_PASSWORD`                      | Redis password (if any)                             |
| `CAMPUS_CART_REDIS_DB`                            | Redis database index                                |
| `CAMPUS_CART_INTEGRATION_RESEND_API_KEY`          | Resend API key for email delivery                   |
| `CAMPUS_CART_AUTH_COOKIE_DOMAIN`                  | Domain for session cookies                          |
| `CAMPUS_CART_CLOUDINARY_CLOUD_NAME`               | Cloudinary cloud name                               |
| `CAMPUS_CART_CLOUDINARY_API_KEY`                  | Cloudinary API key                                  |
| `CAMPUS_CART_CLOUDINARY_API_SECRET`               | Cloudinary API secret                               |
| `CAMPUS_CART_OBSERVABILITY_NEW_RELIC_LICENSE_KEY` | New Relic license key                               |
| `CAMPUS_CART_OBSERVABILITY_NEW_RELIC_APP_NAME`    | New Relic application name                          |

### Client

| Variable              | Description                     |
| --------------------- | ------------------------------- |
| `NEXT_PUBLIC_API_URL` | Backend API base URL (required) |

## Development Workflow

### Server Tasks (Taskfile)

The server uses [Taskfile](https://taskfile.dev) for common development tasks:

| Command                | Description                                     |
| ---------------------- | ----------------------------------------------- |
| `task run`             | Start the Go server (`go run ./cmd/campusCart`) |
| `task tidy`            | Run `go mod tidy`                               |
| `task migrations:new`  | Create a new migration file                     |
| `task migrations:up`   | Apply pending migrations                        |
| `task migrations:down` | Roll back the last migration                    |

### Client Scripts (npm)

| Command         | Description                                         |
| --------------- | --------------------------------------------------- |
| `npm run dev`   | Start the Next.js development server with Turbopack |
| `npm run build` | Build the production bundle                         |
| `npm run start` | Start the production server                         |
| `npm run lint`  | Run ESLint                                          |

### File Upload Utility

The server includes a standalone CLI tool at `cmd/uploader/main.go` for uploading files to Cloudinary via the server's signed upload endpoint. This is useful for seeding media assets or testing uploads outside the browser.

## API Documentation

The server exposes an interactive OpenAPI 3.0.3 documentation page at the `/docs` endpoint. Static files for the spec are served from `apps/server/static/`.

### API Route Summary

**Public Endpoints (No Authentication Required)**

| Method | Path                          | Description                                    |
| ------ | ----------------------------- | ---------------------------------------------- |
| `POST` | `/auth/register`              | Register a new user account                    |
| `POST` | `/auth/login`                 | Log in with email and password                 |
| `POST` | `/auth/verify-email`          | Verify email with 6-digit code                 |
| `GET`  | `/categories`                 | List all categories                            |
| `GET`  | `/categories/{id}`            | Get a specific category                        |
| `GET`  | `/categories/{id}/attributes` | Get category attributes (merged with parent)   |
| `GET`  | `/listings`                   | Browse listings with filters                   |
| `GET`  | `/listings/{id}`              | Get a specific listing (increments view count) |

**Authenticated Endpoints (Session Cookie Required)**

| Method   | Path                         | Description                        |
| -------- | ---------------------------- | ---------------------------------- |
| `POST`   | `/auth/logout`               | Log out (delete session)           |
| `GET`    | `/auth/me`                   | Get currently authenticated user   |
| `GET`    | `/brands/me`                 | Get the current user's brand       |
| `PATCH`  | `/brands/me`                 | Update current user's brand        |
| `POST`   | `/listings`                  | Create a new listing               |
| `PATCH`  | `/listings/{id}`             | Update an existing listing         |
| `DELETE` | `/listings/{id}`             | Delete (soft) a listing            |
| `POST`   | `/listings/upload-signature` | Get a signed Cloudinary upload URL |

**Admin Endpoints (Admin Role Required)**

| Method   | Path               | Description                 |
| -------- | ------------------ | --------------------------- |
| `POST`   | `/categories`      | Create a new category       |
| `PATCH`  | `/categories/{id}` | Update an existing category |
| `DELETE` | `/categories/{id}` | Delete a category           |

### Authentication Flow

1. **Register**: User provides first name, last name, username, email (`@st.ug.edu.gh` required), and password. A verification code is sent to the email via an async background job.
2. **Verify Email**: User submits the 6-digit verification code. On success, a session is created and the `cc_refresh_token` cookie is set (7-day expiry, HttpOnly, SameSite Lax, Secure in production).
3. **Login**: User provides email and password. On success, a new session is created and the cookie is set.
4. **Session Management**: The session uses sliding expiration. Each authenticated request refreshes the session's last activity timestamp and extends the expiry window.
5. **Logout**: The session is deleted server-side and the cookie is cleared.

## Database Schema

The database uses PostgreSQL 16 with the following extensions:

- `uuid-ossp` -- UUID generation
- `citext` -- Case-insensitive text type

### Custom Enums

- `user_role`: admin, user
- `listing_condition`: new, used, refurbished
- `media_type`: image, video
- `feedback_type`: suggestion, bug, other

### Tables

| Table                 | Description                                                                                         |
| --------------------- | --------------------------------------------------------------------------------------------------- |
| `users`               | User accounts with email/phone verification, banning, soft delete                                   |
| `sessions`            | Cookie-based sessions with hashed refresh tokens, IP/UA tracking, sliding expiry                    |
| `categories`          | Hierarchical product categories (parent_id self-reference)                                          |
| `category_attributes` | Dynamic attributes per category (JSONB options, typed: text/number/boolean/enum)                    |
| `brands`              | Seller storefronts with slug, profile/banner images, verification badge                             |
| `listings`            | Product/service listings with JSONB attributes, image/video URL arrays, view counter, promoted flag |
| `saved`               | User bookmarks for listings (unique per user+listing)                                               |
| `reviews`             | 1-5 star ratings with comments and image attachments (unique per listing+reviewer)                  |
| `conversations`       | Chat threads between buyers and sellers                                                             |
| `messages`            | Individual messages within conversations, with media support and read tracking                      |
| `feedback`            | Bug reports, suggestions, and other user feedback with admin notes                                  |

All tables with mutable data include `created_at`, `updated_at` (auto-triggered), and `deleted_at` (soft delete) timestamps.

### Migrations

Migrations are managed using tern and are embedded into the Go binary at compile time. They are located in `apps/server/internal/database/migrations/` and are applied automatically on server startup in non-development environments.

## MVP Features

### Priority 0 (MVP for GA Release)

**Seller Features**

- Account creation with student email verification (`@st.ug.edu.gh`)
- Listing creation with images, videos, condition, description, price, and dynamic category-specific attributes
- Listing update and deletion (soft delete)
- Brand profile management (name, description, profile image, banner image)
- Real-time chat with buyers (planned)
- Bug reporting and feedback submission (planned)
- Subscription payment for listing promotion (planned)
- Account and subscription deletion

**Buyer Features**

- Browse all listings without creating an account
- View individual listings with full details
- Search and filter listings by category, brand, price range, condition, and keyword
- Save/bookmark listings of interest (planned)
- Chat with sellers (planned)
- Rate and review sellers (planned)
- Bug reporting and feedback submission (planned)
- Account deletion

### Priority 1 (Important for Delightful Experience)

- Polished, modern UI/UX across all pages
- Seller analytics dashboard (free for initial period, then subscription-based)
- Advanced search and filtering capabilities
- Fraud reporting for sellers

### Priority 2 (Nice to Have)

- Recommended items based on browsing history and interactions
- Profile image upload for buyers
- Custom background design for seller brand pages

## Post-MVP Roadmap

- **Featured Listings**: Paid promotion to appear at the top of search results (already supported via the `promoted` database flag)
- **Advertising**: Student-friendly brand promotions and student-targeted ads
- **In-App Delivery**: Optional delivery service if the seller opts in, requiring delivery agent infrastructure
- **Enhanced Security**: Additional fraud prevention and detection mechanisms
- **Monetization**: Subscription tiers for sellers once 500+ active sellers are on the platform

## Success Metrics

| Goal                    | Signal                                             | Metric                                       | Target     |
| ----------------------- | -------------------------------------------------- | -------------------------------------------- | ---------- |
| Engagement and Adoption | Sellers and buyers find the product valuable       | User base size                               | Growing    |
| Grow User Base          | Buyers browsing products; Sellers listing products | Weekly active users, 1-week retention rate   | Increasing |
| User Growth Rate        | Consistent month-over-month growth                 | MOM 28-day active user growth rate           | Positive   |
| Business Value          | Sellers and buyers continue using the platform     | Monthly active domains (28-day active users) | Sustained  |

## Contributing

1. Fork the repository.
2. Create a feature branch from `main`.
3. Make your changes following the existing code style and conventions.
4. Ensure all existing functionality continues to work.
5. Submit a pull request with a clear description of the changes.

### Code Style

- **Client**: Follow the ESLint configuration (Next.js core web vitals + TypeScript rules). Use TypeScript strict mode. Follow the shadcn/ui component patterns.
- **Server**: Follow standard Go conventions. Use the existing layered architecture (handler -> service -> repository). All errors should be translated to appropriate HTTP errors using the `err` and `sqlerr` packages.

## License

This project is proprietary. All rights reserved.
