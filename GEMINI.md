# NeoBank Go Project Overview

A microservices-based banking application built with Go, PostgreSQL, and gRPC.

## Architecture

The project follows a microservices architecture with a shared PostgreSQL database using schema-level isolation.

### Core Services

- **Ledger Service (`ledger-service`):**
    - Implements a double-entry accounting system.
    - Manages accounts (`ledger_accounts`), transfers (`transfers`), and entries (`entries`).
    - Exposes a gRPC API for account creation, balance retrieval, and transfer posting.
    - Uses **Serializable** isolation level for ledger transactions to ensure data integrity.
- **Account Service (`account-service`):**
    - Manages customer information and account mapping (WIP).
    - Communicates with the Ledger Service via gRPC.

### Shared Resources

- **`neobank-proto`:** Contains shared Protobuf definitions used for service-to-service communication.
- **PostgreSQL:** Shared database instance with service-specific schemas (`ledger`, `account`).

## Technologies

- **Language:** Go (1.24+)
- **Database:** PostgreSQL (v16)
- **API:** gRPC / Protobuf
- **Logging:** Uber Zap (`go.uber.org/zap`)
- **Migrations:** Liquibase (SQL-based migrations in `db/changelog`)
- **Orchestration:** Docker Compose

## Development Guide

### Building and Running

1. **Infrastructure:** Start the database using Docker Compose.
   ```bash
   docker-compose up -d
   ```
2. **Migrations:** Ensure database schemas are updated (Liquibase). TODO: Add specific command for running migrations.
3. **Run Services:**
   - **Ledger Service:**
     ```bash
     cd ledger-service
     go run main.go
     ```
   - **Account Service:**
     ```bash
     cd account-service
     go run main.go
     ```

### Development Conventions

- **Repository Pattern:** Used for data access (e.g., `ledger_repo.go`).
- **Domain-Driven Design:** Entities and core logic are located in `domain` packages.
- **gRPC API:** Contracts are defined in `.proto` files in `neobank-proto` and generated in service-specific `proto` directories.
- **Logging:** Use `zap.L()` for structured logging throughout the application.
- **Error Handling:** Use wrapped errors for context (e.g., `fmt.Errorf("...: %w", err)`).

### TODOs / Roadmap

- [ ] Complete `account-service` implementation (Customer CRUD, Account mapping).
- [ ] Implement Ledger client in `account-service`.
- [ ] Add unit and integration tests.
- [ ] Finalize Dockerfiles for production deployment.
- [ ] Standardize migration runner script.
