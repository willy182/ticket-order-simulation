# Agent: Ticket Order Simulation

## Overview

This is a Go microservice for simulating a ticket ordering system. It is built using the **Candi** framework (`github.com/golangid/candi`) with a clean hexagonal architecture consisting of two main modules: **Transaction** and **Ticket**.

## Architecture

The project follows a layered architecture:

```
ticketing/
├── main.go                          # Entry point
├── internal/
│   ├── modules/
│   │   ├── ticket/                  # Ticket module (CRUD + cache)
│   │   │   ├── delivery/
│   │   │   │   └── resthandler/     # REST API handlers
│   │   │   ├── domain/              # Domain models, filters, payloads
│   │   │   ├── repository/          # Database repository interfaces & impl
│   │   │   ├── usecase/             # Business logic
│   │   │   └── module.go            # Module wiring
│   │   └── transaction/             # Transaction module (order processing)
│   │       ├── delivery/
│   │       │   ├── resthandler/     # REST API handlers
│   │       │   └── workerhandler/   # Task queue (background job) handlers
│   │       ├── domain/              # Domain models, filters, payloads
│   │       ├── repository/          # Database repository interfaces & impl
│   │       ├── usecase/             # Business logic + Lua script for stock
│   │       └── module.go            # Module wiring
│   └── pkg/
│       ├── helper/                  # Shared helpers (JWT, pagination, constants)
│       ├── shared/
│       │   ├── domain/              # Shared domain models (Transaction, Ticket)
│       │   ├── repository/          # Shared DB repository (PostgreSQL + MongoDB)
│       │   └── usecase/             # Shared usecase interface & wiring
│       └── infra/                   # Infrastructure (Redis cache, etc.)
├── api/
│   ├── jsonschema/                  # JSON schema validation files
│   │   └── transaction/save.json
│   └── docs/                        # Swagger/OpenAPI docs (if present)
├── configs/
│   └── configs.go                   # Configuration loader
└── migration/                       # Database migrations
```

## Technology Stack

| Component       | Technology                          |
|-----------------|-------------------------------------|
| Language        | Go (Golang)                         |
| Framework       | Candi v1.20.0                       |
| Database        | PostgreSQL (GORM) + MongoDB        |
| Cache           | Redis                               |
| Task Queue      | Built-in Candi task queue worker    |
| Validation      | JSON Schema + Candi validator       |
| Tracing         | Candi tracer (OpenTelemetry-compatible) |
| Dependency Injection | Candi dependency framework      |

## Modules

### 1. Ticket Module

Manages ticket inventory — CRUD operations with Redis caching.

**Domain Models** (`internal/modules/ticket/domain/`):
- `RequestTicket` — input payload for creating/updating tickets
- `ResponseTicket` — serialized output model
- `FilterTicket` — query filter with pagination
- `PayloadTicket` — MongoDB-specific payload

**Use Cases** (`internal/modules/ticket/usecase/`):
- `CreateTicket` — creates a ticket and caches it in Redis under key `kuota_ticket_{id}`
- `GetAllTicket` — fetches all tickets with pagination; uses `errgroup` for concurrent count/fetch
- `GetDetailTicket` — fetches a single ticket by ID
- `UpdateTicket` — updates ticket fields (Title, Quota, Price) and updates Redis cache
- `DeleteTicket` — deletes a ticket by ID

**Delivery** (`internal/modules/ticket/delivery/resthandler/`):
- REST endpoints under `/v1/ticket/` with Swagger documentation
- JSON schema validation via `ticket/save`

### 2. Transaction Module

Handles ticket ordering with async processing via task queues.

**Domain Models** (`internal/modules/transaction/domain/`):
- `RequestTransaction` — input payload (customer info, ticket ID, quantity, status)
- `ResponseTransaction` — serialized output model
- `FilterTransaction` — query filter with date range support
- `ReqSendEmail` — email notification payload

**Use Cases** (`internal/modules/transaction/usecase/`):
- `CreateTransaction` — enqueues a `ReserveTicket` task to the task queue (async)
- `ReserveTicket` — the core business logic executed by the worker:
  1. Reads ticket data from Redis cache
  2. Uses a **Lua script** for atomic stock deduction in Redis (handles stock not found / sold out)
  3. Generates an invoice number and calculates total amount
  4. Saves to PostgreSQL
  5. Enqueues `GenerateTicketCode` task
  6. On failure status, enqueues `SendEmail` task
- `GenerateTicketCode` — generates a unique ticket code and enqueues `SendEmail`
- `SendEmail` — simulates sending an email notification (currently returns a formatted error message)
- `GetAllTransaction` / `GetDetailTransaction` / `UpdateTransaction` / `DeleteTransaction` — standard CRUD

**Delivery** (`internal/modules/transaction/delivery/`):
- `resthandler/` — REST endpoints under `/v1/transaction/` 
- `workerhandler/` — Task queue handlers for `ReserveTicket`, `SendEmail`, `GenerateTicketCode`

## Request Flow

```
1. Client POST /v1/transaction/
2. REST handler validates payload via JSON Schema
3. CreateTransaction usecase enqueues "ReserveTicket" task
4. Task queue worker picks up the job
5. ReserveTicket:
   a. Read ticket from Redis (kuota_ticket_{id})
   b. Execute Lua script to atomically decrement stock
   c. If stock insufficient → return error (auto-retry up to 3x)
   d. If success → save to PostgreSQL, enqueue GenerateTicketCode
6. GenerateTicketCode:
   a. Generate random ticket code
   b. Update PostgreSQL
   c. Enqueue SendEmail
7. SendEmail:
   a. Format and send notification email
```

## Redis Lua Script

The `ReserveTicket` usecase uses an embedded Lua script for atomic stock operations:
- Returns `-2` if the key does not exist (stock not found)
- Returns `-1` if stock is insufficient (sold out)
- Returns the new stock count on success

## JSON Schema Validation

Schemas are defined in `api/jsonschema/` and loaded by the Candi validator. The `transaction/save` schema validates both create and update requests.

## Configuration

Managed via `configs/configs.go` using the Candi config system with environment variable support and `.env` file loading.

## Project Generation

This project was scaffolded using **Candi v1.20.0** code generator (`github.com/golangid/candi`).