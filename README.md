# Document Processing Orchestrator

A production-inspired distributed document processing orchestrator built in **Go**.

The project explores reliable asynchronous workflow execution using database-backed jobs, worker coordination, job leasing, retries, heartbeats and failure recovery.

The main goal is to process documents through multiple asynchronous stages while keeping the system observable, resilient and horizontally scalable.

---

## Overview

The orchestrator receives documents and coordinates their processing through a sequence of jobs.

Example workflow:

```text
Upload Document
      ↓
Validate Document
      ↓
Extract Content
      ↓
Extract Metadata
      ↓
Enrich Data
      ↓
Build Result
      ↓
Completed
```

Each stage is executed asynchronously by workers.

The orchestration layer is responsible for:

- job scheduling;
- job state transitions;
- worker coordination;
- retries;
- failure recovery;
- job leasing;
- heartbeats;
- workflow dependencies;
- observability.

---

## Current Architecture

```text
                    ┌─────────────────┐
                    │     Client      │
                    └────────┬────────┘
                             │
                             │ HTTP
                             ▼
                    ┌─────────────────┐
                    │     Go API      │
                    └────────┬────────┘
                             │
                             │ database/sql
                             ▼
                    ┌─────────────────┐
                    │     MariaDB     │
                    │                 │
                    │ Jobs + State    │
                    └───────┬─────────┘
                            │
                    acquire │ jobs
                            │
                ┌───────────┴───────────┐
                ▼                       ▼
         ┌─────────────┐         ┌─────────────┐
         │   Runner    │         │   Runner    │
         │      A      │         │      B      │
         └─────────────┘         └─────────────┘
```

At the current stage, the project contains:

- Go API;
- Go runner;
- MariaDB;
- automatic database migrations;
- Docker Compose development environment;
- structured logging;
- health endpoint;
- initial jobs persistence model.

---

## Tech Stack

### Backend

- Go
- `net/http`
- Chi Router
- `database/sql`

### Database

- MariaDB
- SQL migrations
- `golang-migrate`

### Infrastructure

- Docker
- Docker Compose
- Multi-stage Docker builds
- Makefile

### Observability

- Structured logging with `log/slog`

Additional observability features will be added as the project evolves.

---

## Project Structure

```text
.
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── runner/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   └── database/
│       └── mysql.go
│
├── migrations/
│   ├── 000001_create_jobs.up.sql
│   └── 000001_create_jobs.down.sql
│
├── docs/
│   ├── architecture.md
│   └── decisions/
│
├── Dockerfile.api
├── Dockerfile.runner
├── docker-compose.yml
├── Makefile
├── .env.example
├── go.mod
└── README.md
```

The project structure will evolve as document workflows, job orchestration and worker components are introduced.

---

## Running Locally

### Requirements

You only need:

- Docker
- Docker Compose
- Make

Go is only required if you want to run the applications directly outside Docker.

---

## Environment Configuration

Create your local environment file:

```bash
cp .env.example .env
```

Example:

```env
APP_ENV=development

HTTP_PORT=8080

DB_HOST=mariadb
DB_PORT=3306
DB_NAME=orchestrator
DB_USER=orchestrator
DB_PASSWORD=orchestrator

DB_ROOT_PASSWORD=root
```

---

## Start the Environment

Run:

```bash
make up
```

Or directly:

```bash
docker compose up -d --build
```

This starts:

```text
MariaDB
   ↓
Migrations
   ↓
API
   ↓
Runner
```

Database migrations are executed automatically before the API and runner start.

---

## Check Running Services

```bash
make ps
```

Or:

```bash
docker compose ps
```

Expected services:

```text
orchestrator-mariadb
orchestrator-api
orchestrator-runner
```

The migration container is a one-shot service and exits after applying pending migrations.

---

## Health Check

The API exposes:

```http
GET /health
```

Example:

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "ok"
}
```

---

## Database Migrations

Database schema changes are managed using `golang-migrate`.

### Create a migration

```bash
make migration name=create_documents
```

This generates:

```text
migrations/
├── 000002_create_documents.up.sql
└── 000002_create_documents.down.sql
```

### Apply migrations

```bash
make migrate-up
```

### Roll back the latest migration

```bash
make migrate-down
```

### Check current migration version

```bash
make migrate-version
```

---

## Development Commands

### Start

```bash
make up
```

### Stop

```bash
make down
```

### Build containers

```bash
make build
```

### Follow logs

```bash
make logs
```

### Run tests

```bash
make test
```

### Update Go modules

```bash
make tidy
```

### Reset the development database

```bash
make db-reset
```

This removes the local database volume and recreates the environment from scratch.

---

## Job Model

Jobs are persisted in MariaDB.

The initial job model contains information required for future distributed execution:

```text
id
type
payload
status

attempts
max_attempts

lease_owner
lease_expires_at
last_heartbeat_at

available_at
last_error

created_at
updated_at
```

These fields will allow the orchestrator to implement:

- worker leasing;
- retry policies;
- abandoned job recovery;
- scheduled execution;
- job ownership;
- horizontal worker scaling.

---

## Planned Workflow

The initial document processing workflow will evolve toward:

```text
DOCUMENT_CREATED
      ↓
VALIDATE_DOCUMENT
      ↓
EXTRACT_CONTENT
      ↓
EXTRACT_METADATA
      ↓
ENRICH_DATA
      ↓
BUILD_RESULT
      ↓
COMPLETED
```

Some stages may later execute in parallel.

Example:

```text
                  ┌── EXTRACT_METADATA ──┐
EXTRACT_CONTENT ──┤                      ├── BUILD_RESULT
                  └── GENERATE_PREVIEW ─┘
```

This will introduce workflow dependencies and DAG-based orchestration.

---

## Roadmap

### Phase 1 — Foundation

- [x] Go project structure
- [x] Docker Compose environment
- [x] MariaDB
- [x] Automatic database migrations
- [x] API service
- [x] Runner service
- [x] Structured logging
- [x] Health endpoint
- [x] Initial jobs table

### Phase 2 — Document Processing

- [ ] Document persistence
- [ ] `POST /documents`
- [ ] Initial document validation job
- [ ] Job acquisition by workers
- [ ] Job execution
- [ ] Job completion and failure states

### Phase 3 — Reliability

- [ ] Transactional job leasing
- [ ] Worker heartbeats
- [ ] Expired lease recovery
- [ ] Retry strategies
- [ ] Exponential backoff
- [ ] Maximum attempts
- [ ] Dead-letter / permanent failure handling
- [ ] Graceful shutdown

### Phase 4 — Workflow Orchestration

- [ ] Multiple processing stages
- [ ] Job dependencies
- [ ] Parallel jobs
- [ ] Workflow state management
- [ ] Idempotency
- [ ] Job cancellation
- [ ] Scheduled jobs

### Phase 5 — Observability

- [ ] Prometheus metrics
- [ ] Job processing metrics
- [ ] Worker metrics
- [ ] HTTP metrics
- [ ] Request IDs
- [ ] Distributed tracing
- [ ] Readiness checks

### Phase 6 — Engineering Quality

- [ ] Unit tests
- [ ] Integration tests
- [ ] Concurrent worker tests
- [ ] Race detector
- [ ] GitHub Actions
- [ ] Linting
- [ ] Architecture Decision Records

---

## Design Principles

This project follows a few core principles.

### Reliability over convenience

Jobs should remain recoverable even when workers crash or infrastructure fails.

### Explicit state transitions

Workflow state should be persisted and observable instead of living only in process memory.

### Horizontal execution

Multiple workers should be able to safely compete for work.

### Idempotent processing

Retrying a job should not corrupt the workflow or generate duplicate side effects.

### Simple infrastructure

The initial design intentionally uses MariaDB as both persistence and coordination storage to reduce operational complexity while exploring transactional job orchestration.

---

## Architecture Decisions

Important design decisions are documented as ADRs under:

```text
docs/decisions/
```

Example:

```text
0001-database-backed-queue.md
```

The goal is to document not only **what** the system does, but also **why** architectural decisions were made and which trade-offs were considered.

---

## Why MariaDB as a Job Queue?

The first version intentionally uses the relational database for job coordination instead of introducing Redis, RabbitMQ or Kafka.

This allows the project to explore:

- transactional job acquisition;
- row-level locking;
- durable state;
- worker coordination;
- job recovery;
- consistency between workflow state and job state.

As the system evolves, alternative messaging architectures may be evaluated and documented.

---

## Status

🚧 **Work in progress**

The current version provides the infrastructure foundation for the orchestrator.

The next milestone is implementing document creation and the first real processing workflow:

```text
POST /documents
      ↓
create document
      ↓
create VALIDATE_DOCUMENT job
      ↓
runner acquires job
      ↓
process document
      ↓
complete job
```

---

## Author

**José Vitor Alves Pignataro**

Senior Software Engineer & Tech Lead focused on backend architecture, distributed systems and scalable platforms.

- GitHub: [@Vitorpignataro](https://github.com/Vitorpignataro)
- LinkedIn: [José Vitor Alves Pignataro](https://www.linkedin.com/in/jos%C3%A9-vitor-alves-pignataro/)
