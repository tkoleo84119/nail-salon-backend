# CLAUDE.md

This file provides guidance to Claude Code when working with this nail salon backend codebase.

## Development Commands

- `make run` - Run the application server
- `make test` - Run all tests
- `make seed-test` - Seed test data
- `make migrate-up` - Run database migrations up
- `make migrate-down` - Run database migrations down

### Go Commands
- `go run cmd/server/main.go` - Direct server execution
- `go test ./...` - Run all tests

## Architecture

This is a Go web API for a nail salon management system using:
- **Framework**: Gin HTTP framework (`github.com/gin-gonic/gin`)
- **Database**: PostgreSQL with pgx/v5 (`github.com/jackc/pgx/v5`) and sqlx (`github.com/jmoiron/sqlx`)
- **Code Generation**: SQLC for type-safe database queries
- **Authentication**: JWT tokens (`github.com/golang-jwt/jwt/v5`) with role-based access control
- **ID Generation**: Snowflake algorithm (`github.com/bwmarrin/snowflake`) for unique IDs
- **Password Hashing**: bcrypt (`golang.org/x/crypto/bcrypt`)
- **Configuration**: Environment-based config with .env file support (`github.com/joho/godotenv`)

## Project Structure

```
nail-salon-backend/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── app/            # Application container and routing
│   ├── config/         # Configuration management
│   ├── errors/         # Centralized error handling
│   ├── handler/        # HTTP handlers (controllers)
│   ├── infra/          # Infrastructure layer (database)
│   ├── middleware/     # HTTP middleware
│   ├── model/          # DTOs and data models
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   └── utils/          # Shared utilities
├── migration/          # Database migrations
├── scripts/           # Utility scripts
└── docs/              # API documentation
```

## Universal Implementation Pattern

All APIs (admin and customer) follow the same 3-layer pattern:

### File Organization
```
{context}/             # 'admin' or root level
├── {domain}/          # auth, booking, store, etc.
│   ├── {operation}.go # create, get, update, etc.
│   └── interface.go   # service interfaces (in service layer)
```

### Naming Conventions
- **Files**: `{operation}.go` (e.g., `create.go`, `get_all.go`, `update_my.go`)
- **Structs**: `type {Operation} struct` (e.g., `type Create struct`)
- **Models**: `{Operation}Request/Response` (e.g., `CreateRequest`)
- **Package Aliases**: `{context}{Domain}Handler/Service/Model`

### Layer Pattern

**1. Handler Layer**
```go
type {Operation} struct {
    service {domain}Service.{Operation}Interface
}

func (h *{Operation}) {Operation}(c *gin.Context) {
    // 1. Bind & validate request
    var req {domain}Model.{Operation}Request
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErrors := utils.ExtractValidationErrors(err)
        errorCodes.RespondWithValidationErrors(c, validationErrors)
        return
    }

    // 2. Extract auth context (if protected)
    userContext, exists := middleware.Get{User}FromContext(c)

    // 3. Parse string IDs to int64 (if needed)
    parsedID, err := utils.ParseID(req.ID)

    // 4. Call service
    response, err := h.service.{Operation}(c.Request.Context(), req, userContext)
    if err != nil {
        errorCodes.RespondWithServiceError(c, err)
        return
    }

    // 5. Return response
    c.JSON(http.StatusOK, common.SuccessResponse(response))
}
```

**2. Service Layer**
```go
// Interface (in interface.go)
type {Operation}Interface interface {
    {Operation}(ctx context.Context, req Model.{Operation}Request, ...) (*Model.{Operation}Response, error)
}

// Implementation
type {Operation} struct {
    queries *dbgen.Queries      // SQLC for simple ops
    db      *pgxpool.Pool       // Transactions
    repo    Repository          // Complex queries
}

func (s *{Operation}) {Operation}(ctx context.Context, req Model.{Operation}Request, ...) (*Model.{Operation}Response, error) {
    // 1. Validate business rules & permissions
    // 2. Database operations (with tx if needed)
    // 3. Return formatted response
}
```

**3. Model Layer**
```go
type {Operation}Request struct {
    Field1 string   `json:"field1" binding:"required,max=50"`
    IDs    []string `json:"ids" binding:"required,min=1,max=10"`
}

type {Operation}Response struct {
    ID        string `json:"id"`
    CreatedAt string `json:"createdAt"`
}
```

## Database Usage Guide

- **SQLC (queries)**: Standard CRUD, single-table operations, type-safe queries
- **SQLX (repo)**: Dynamic queries, complex joins, partial updates, filtering
- **PGX Pool (db)**: Transactions, bulk operations, high-performance needs

## Authentication

**Customer Auth**: LINE OAuth → `GetCustomerFromContext(c)`
**Staff Auth**: Username/password → `GetStaffFromContext(c)`

**Staff Roles** (hierarchical): `SUPER_ADMIN` > `ADMIN` > `MANAGER` > `STYLIST`

## Error Handling

Service errors use hierarchical codes: `AUTH_*`, `BOOKING_*`, `CUSTOMER_*`, etc.
```go
return errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
errorCodes.RespondWithServiceError(c, err)
```

## Common Utilities

**ID Handling**:
```go
newID := utils.GenerateID()                    // Generate Snowflake ID
parsed := utils.ParseID(stringID)              // String → int64
formatted := utils.FormatID(int64ID)           // int64 → string
```

**Time Handling**:
```go
pgTime := utils.TimeToPgTimestamptz(time.Now())
timeString := utils.PgTimestamptzToTimeString(pgTime)
```

## AI Implementation Steps

### Adding New Feature
1. **Identify domain & context** (admin vs customer)
2. **Create model** in `/internal/model/{context}/{domain}/{operation}.go`
3. **Define interface** in `/internal/service/{context}/{domain}/interface.go`
4. **Implement service** in `/internal/service/{context}/{domain}/{operation}.go`
5. **Create handler** in `/internal/handler/{context}/{domain}/{operation}.go`
6. **Register in container** (`container_{context}.go`)
7. **Add route** in router setup

### Standard Operations
- `create.go`, `get.go`, `get_all.go`, `update.go`, `delete.go`
- `get_my.go`, `update_my.go` (user's own data)
- `{operation}_bulk.go` (bulk operations)

### Key Patterns
- Always use interface-based dependency injection
- Validate at both request and business logic levels
- Use transactions for multi-step operations
- Extract auth context for protected endpoints
- Return structured error responses
- Convert string IDs to int64 for database operations

## Environment Variables
- `DB_URL`, `JWT_SECRET`, `LINE_CHANNEL_ID`, `LINE_CHANNEL_SECRET`, `PORT`