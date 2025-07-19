# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Core Commands
- `make run` - Run the application server
- `make test` - Run all tests
- `make seed-test` - Seed test data into database
- `make migrate-up` - Run database migrations up (requires DB_URL env var)
- `make migrate-down` - Run database migrations down (requires DB_URL env var and NUMBER)

### Go Commands
- `go run cmd/server/main.go` - Direct server execution
- `go test ./...` - Run all tests
- `go run scripts/seed/test_seed.go` - Seed test data

## Architecture Overview

This is a Go web API for a nail salon management system using:
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with both pgx/v5 and sqlx drivers
- **Code Generation**: SQLC for type-safe database queries
- **Authentication**: JWT tokens with role-based access control
- **ID Generation**: Snowflake algorithm for unique IDs
- **Configuration**: Environment-based config with .env file support

### Project Structure
```
internal/
├── config/          # Environment configuration
├── errors/          # Centralized error management with YAML definitions
├── handler/         # HTTP handlers (auth, staff, store-access)
├── infra/db/        # Database connection setup
├── middleware/      # JWT auth and role authorization
├── model/           # Request/response models
├── repository/      # Data access layer
│   ├── sqlc/        # SQLC generated queries
│   └── sqlx/        # SQLX manual queries
├── service/         # Business logic layer
├── testutils/       # Test utilities and mocks
└── utils/           # Utility functions (JWT, passwords, validation)
```

### Key Patterns
- **Layered Architecture**: Handler → Service → Repository
- **Dependency Injection**: Services injected into handlers via constructors
- **Error Handling**: Centralized error codes in `internal/errors/errors.yaml`
- **Authentication**: JWT middleware with role-based permissions
- **Database**: Uses both SQLC (for type safety) and SQLX (for flexibility)
- **Testing**: Comprehensive unit tests with mocks for external dependencies

### Database
- Uses PostgreSQL with migrations in `migration/` directory
- SQLC configuration in `sqlc.yaml` generates type-safe queries
- Schema files in @docs/db/database.dbml

### Environment Configuration
Required environment variables:
- `JWT_SECRET` - JWT signing secret
- Database connection: `DB_DSN` or individual DB_* variables
- Optional: `PORT`, `SNOWFLAKE_NODE_ID`, JWT/DB timing configs

### Business Domain
Nail salon management system with:
- Staff management with role-based access (admin/staff)
- Store access permissions for multi-location support
- Authentication and authorization
- Future: customer management, appointments, inventory, billing

### Commit Convention
Follow Conventional Commits format: `<type>: <description>`
- Types: feat, fix, refactor, perf, style, test, docs, build, ops, chore
- Use imperative mood, present tense, no capitalization, no period
- English language for descriptions

### Testing Guidelines
- For handler & service, always write corresponding test files
- All test files need to passed before commit
- Only write test when @internal/model/ files have function
- If service use sqlx, the mock repository need to write in @internal/testutils/mocks/repository.go 

### Development Notes
- If needs new SQL query, look up current SQL query first (in `internal/repository/sqlc/` and `internal/repository/sqlx/`), not create duplicate SQL
- If sql is update optional columns(dynamic), use sqlx not sqlc
- When adding a new SQL for batch insert operations, use sqlc's :copyfrom
- Handler & service & model need to use modules struct by business domain (separate by folder)

### Validation Patterns

#### Handler Layer Validation (Fixed Order)
1. **Input JSON Validation** - Always first, using `c.ShouldBindJSON(&req)`
2. **Path Parameter Validation** - Validate required URL parameters
3. **Business Logic Validation** - Check if update requests have fields using `req.HasUpdates()` or `req.HasUpdate()`
4. **Authentication Context Validation** - Use `middleware.GetStaffFromContext(c)` (not legacy `c.Get("staffContext")`)
5. **ID Parsing Validation** - Convert string IDs using `utils.ParseID()`
6. **Service Layer Call** - Pass validated data to service

#### Service Layer Validation (Fixed Order)
1. **Input Data Validation** - Parse IDs using `utils.ParseID()` and `utils.ParseIDSlice()`
2. **Request Completeness** - Validate update requests have at least one field
3. **Business Logic Validation** - Role validation, time ranges, entity existence
4. **Permission & Authorization** - Role-based checks, store access, ownership validation
5. **Data Integrity Validation** - Uniqueness, conflict prevention, entity state

#### Validation Guidelines
- **Handler Validation**: Focus on input format, authentication context, and basic parameter validation
- **Service Validation**: Handle business logic, permissions, and data integrity
- **Error Handling**: Use `errorCodes.AbortWithError()` in handlers, `errorCodes.NewServiceError()` in services
- **Context Extraction**: Always use `middleware.GetStaffFromContext(c)` for staff context
- **ID Parsing**: Use `utils.ParseID()` for string to int64 conversion with validation
- **Time Validation**: Use `common.ParseTimeSlot()` for time format validation
- **Update Validation**: Implement `HasUpdates()` or `HasUpdate()` methods on request models

### Error Handling Patterns
- Use predefined error constants for consistent error management
- For permission issues, use `AUTH.AUTH_PERMISSION_DENIED` instead of creating new error