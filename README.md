# Agentcore Backend Service

A production-ready Go backend service with Gin framework.

## Features

- RESTful API endpoints
- PostgreSQL database integration
- Environment configuration
- Logging
- Testing
- Linting support

## Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Make

## Getting Started

1. Clone the repository
2. Install dependencies:

```bash
make install-deps
```

3. Set up the database by running the schema:

```bash
psql -U your_user -d your_database -f db/schema.sql
```

4. Create a `.env` file based on the example:

```bash
cp .env.example .env
```

5. Run the application:

```bash
make run
```

## Commands

- `make run` - Start the application
- `make test` - Run tests
- `make lint` - Run linter
- `make build` - Build the application
- `make test-coverage` - Run tests with coverage
- `make docs` - Generate documentation

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `POST /api/v1/query` - Query endpoint with prompt + context metadata

OpenAPI specification is available at `openapi/openapi.yaml`.

## Directory Structure

```
.
├── main.go             # Entry point
├── go.mod              # Go module file
├── Makefile            # Build commands
├── README.md           # This file
├── .env.example        # Environment variables example
├── db/
│   ├── database.go     # Database connection logic
│   └── schema.sql      # Database schema
├── config/
│   └── config.go       # Configuration handling
├── server/
│   └── server.go       # Server setup and routing
├── handlers/
│   └── user.go         # HTTP handlers
├── models/
│   └── user.go         # Data models
└── services/
    └── user.go         # Business logic
```