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
- **Framework**: Gin HTTP framework (`github.com/gin-gonic/gin`)
- **Database**: PostgreSQL with pgx/v5 (`github.com/jackc/pgx/v5`) and sqlx (`github.com/jmoiron/sqlx`)
- **Code Generation**: SQLC for type-safe database queries
- **Authentication**: JWT tokens (`github.com/golang-jwt/jwt/v5`) with role-based access control
- **ID Generation**: Snowflake algorithm (`github.com/bwmarrin/snowflake`) for unique IDs
- **Password Hashing**: bcrypt (`golang.org/x/crypto/bcrypt`)
- **Configuration**: Environment-based config with .env file support (`github.com/joho/godotenv`)
- **Testing**: Testify framework (`github.com/stretchr/testify`)

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

### Database Architecture
- **PostgreSQL**: Primary database with connection pooling via `pgxpool.Pool` 
- **pgx/v5**: Primary driver for transactions and connection management
- **SQLC**: Type-safe SQL query generation for standard CRUD operations
- **SQLX**: Dynamic SQL queries for complex updates with optional fields
- **Migrations**: Located in `migration/` directory
- **Schema**: Documented in `docs/db/database.dbml`

### Database Usage Patterns
- **Create/Read/Simple Updates**: Use SQLC for type safety and performance
- **Dynamic Updates**: Use SQLX for optional field updates (e.g., `UpdateStaffRequest`)
- **Batch Operations**: Use SQLC's `:copyfrom` syntax for efficient bulk inserts
- **Transactions**: Always wrap multi-step operations in `pgx.Tx` transactions

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
- **Centralized Error Management**: Error codes defined in `internal/errors/errors.yaml`
- **Service Errors**: `errorCodes.NewServiceError()` and `errorCodes.NewServiceErrorWithCode()`
- **Handler Errors**: `errorCodes.AbortWithError()` and `errorCodes.RespondWithServiceError()`
- **Error Categories**: AUTH, USER, VAL (validation), SYS (system), plus business domain errors
- For permission issues, use `AUTH.AUTH_PERMISSION_DENIED` instead of creating new error

## CRUD Operation Patterns

### CREATE Operations

#### Handler Pattern (Create)
```go
func (h *Handler) CreateEntity(c *gin.Context) {
    // 1. Input JSON validation (always first)
    var req EntityCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErrors := utils.ExtractValidationErrors(err)
        errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
        return
    }

    // 2. Authentication context validation
    staffContext, exists := middleware.GetStaffFromContext(c)
    if !exists {
        errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
        return
    }

    // 3. Service layer call
    response, err := h.service.CreateEntity(c.Request.Context(), req, *staffContext)
    if err != nil {
        errorCodes.RespondWithServiceError(c, err)
        return
    }

    // 4. Success response (201 Created)
    c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
```

#### Service Pattern (Create) 
```go
func (s *Service) CreateEntity(ctx context.Context, req EntityRequest, staffContext StaffContext) (*EntityResponse, error) {
    // 1. Input validation & ID parsing
    entityID, err := utils.ParseID(req.ID)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid ID", err)
    }

    // 2. Business logic validation (permissions, role)
    if err := s.validatePermissions(staffContext.Role, req); err != nil {
        return nil, err
    }

    // 3. Data integrity validation (existence, uniqueness)
    exists, err := s.queries.CheckEntityExists(ctx, req.Name)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "check failed", err)
    }
    if exists {
        return nil, errorCodes.NewServiceErrorWithCode(errorCodes.EntityAlreadyExists)
    }

    // 4. Transaction-based creation
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "transaction failed", err)
    }
    defer tx.Rollback(ctx)

    qtx := dbgen.New(tx)
    created, err := qtx.CreateEntity(ctx, params)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "create failed", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "commit failed", err)
    }

    return buildResponse(created), nil
}
```

#### Repository Pattern (Create)
- **SQLC**: Use for standard single-record creation with `INSERT ... RETURNING`
- **Batch Creates**: Use `:copyfrom` syntax for efficient bulk inserts
- **Transactions**: All creates are wrapped in transactions for atomicity

### READ Operations

#### Query Patterns
- **Single Record**: `GetEntityByID :one`
- **Multiple Records**: `GetEntitiesByIDs :many` 
- **List Operations**: `GetAllActiveEntities :many`
- **Existence Checks**: `CheckEntityExists :one` returning boolean
- **Complex Joins**: SQLC queries with relationships

#### Service Layer (Read)
- Minimal business logic, mainly permission checks
- Direct SQLC query calls (no transactions needed)
- Role-based filtering and store access validation

### UPDATE Operations

#### Handler Pattern (Update)
```go
func (h *Handler) UpdateEntity(c *gin.Context) {
    // 1. Path parameter validation
    targetID := c.Param("id")
    if targetID == "" {
        errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
            "id": "id為必填項目",
        })
        return
    }

    // 2. Input JSON validation
    var req EntityUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErrors := utils.ExtractValidationErrors(err)
        errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
        return
    }

    // 3. Business logic validation - HasUpdates check
    if !req.HasUpdates() {
        errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, map[string]string{
            "request": "至少需要提供一個欄位進行更新",
        })
        return
    }

    // 4. Authentication context validation
    staffContext, exists := middleware.GetStaffFromContext(c)
    if !exists {
        errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
        return
    }

    // 5. Service call
    response, err := h.service.UpdateEntity(c.Request.Context(), targetID, req, staffContext)
    if err != nil {
        errorCodes.RespondWithServiceError(c, err)
        return
    }

    c.JSON(http.StatusOK, common.SuccessResponse(response))
}
```

#### Update Repository Pattern (SQLX for Dynamic Updates)
```go
func (r *Repository) UpdateEntity(ctx context.Context, id int64, req UpdateRequest) (*UpdateResponse, error) {
    setParts := []string{"updated_at = NOW()"}
    args := map[string]interface{}{"id": id}

    // Dynamic field updates
    if req.Field1 != nil {
        setParts = append(setParts, "field1 = :field1")
        args["field1"] = *req.Field1
    }

    if req.Field2 != nil {
        setParts = append(setParts, "field2 = :field2") 
        args["field2"] = *req.Field2
    }

    query := fmt.Sprintf(`
        UPDATE entities SET %s WHERE id = :id
        RETURNING id, field1, field2, updated_at
    `, strings.Join(setParts, ", "))

    var result EntityModel
    rows, err := r.db.NamedQuery(query, args)
    if err != nil {
        return nil, fmt.Errorf("update failed: %w", err)
    }
    defer rows.Close()

    if !rows.Next() {
        return nil, fmt.Errorf("no rows returned")
    }

    if err := rows.StructScan(&result); err != nil {
        return nil, fmt.Errorf("scan failed: %w", err)
    }

    return buildResponse(result), nil
}
```

### DELETE Operations

#### Delete Patterns
- **Soft Delete**: Update `is_active = false` flag (preferred)
- **Hard Delete**: Actual record removal (use cautiously)
- **Bulk Delete**: Multiple records in single transaction
- **Cascade Delete**: Database-level foreign key constraints handle related records

#### Service Pattern (Delete)
```go
func (s *Service) DeleteEntity(ctx context.Context, entityID string, staffContext StaffContext) (*DeleteResponse, error) {
    // 1. ID parsing and validation
    id, err := utils.ParseID(entityID)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid ID", err)
    }

    // 2. Entity existence check
    entity, err := s.queries.GetEntityByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errorCodes.NewServiceErrorWithCode(errorCodes.EntityNotFound)
        }
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "get failed", err)
    }

    // 3. Permission validation
    if err := s.validateDeletePermissions(staffContext, entity); err != nil {
        return nil, err
    }

    // 4. Business constraint validation
    hasConstraints, err := s.queries.CheckEntityConstraints(ctx, id)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "constraint check failed", err)
    }
    if hasConstraints {
        return nil, errorCodes.NewServiceErrorWithCode(errorCodes.EntityConstraintViolation)
    }

    // 5. Perform delete (usually soft delete)
    if err := s.queries.DeleteEntity(ctx, id); err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "delete failed", err)
    }

    return buildDeleteResponse(entity), nil
}
```

## Model Validation Patterns

### Request Model Structure
```go
// Update requests use pointer types for optional fields
type UpdateEntityRequest struct {
    Field1   *string `json:"field1,omitempty" binding:"omitempty,max=100"`
    Field2   *bool   `json:"field2,omitempty"`
    Field3   *int    `json:"field3,omitempty" binding:"omitempty,min=1"`
}

// HasUpdates method required for all update requests
func (r UpdateEntityRequest) HasUpdates() bool {
    return r.Field1 != nil || r.Field2 != nil || r.Field3 != nil
}
```

### Validation Tags
- `binding:"required"` - Required field validation
- `binding:"omitempty"` - Optional field validation (skip if not provided)
- `binding:"oneof=VALUE1 VALUE2"` - Enum validation
- `binding:"email"` - Email format validation
- `binding:"min=X,max=Y"` - Length/value constraints

## Authentication & Authorization Patterns

### JWT Middleware Usage
```go
// Staff authentication (most endpoints)
staffContext, exists := middleware.GetStaffFromContext(c)
if !exists {
    errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
    return
}

// Customer authentication (customer endpoints)  
customerContext, exists := middleware.GetCustomerFromContext(c)
if !exists {
    errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
    return
}
```

### Role-Based Authorization
- **SUPER_ADMIN**: Full system access, can manage all entities
- **ADMIN**: Store-level administration, can manage assigned stores
- **MANAGER**: Store operations, limited administrative functions
- **STYLIST**: Self-service only, can only manage own records

### Permission Validation Patterns
- **Store Access**: Users can only access stores in their `StoreList`
- **Self-Service**: STYLIST role limited to own records (match `staffContext.UserID`)
- **Administrative**: SUPER_ADMIN and ADMIN can manage other users' records

## Transaction Management

### Standard Transaction Pattern
```go
tx, err := s.db.Begin(ctx)
if err != nil {
    return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "transaction failed", err)
}
defer tx.Rollback(ctx) // Always defer rollback for cleanup

qtx := dbgen.New(tx) // Create SQLC queries with transaction
// Perform database operations...

if err := tx.Commit(ctx); err != nil {
    return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "commit failed", err)
}
```

### When to Use Transactions
- **Multi-table operations**: Creating related records across tables
- **Batch operations**: Multiple inserts/updates that must succeed together
- **Business consistency**: Operations that must be atomic for data integrity
- **Single operations**: Simple creates/updates don't need explicit transactions

## Key Development Conventions

### ID Management
- **Generation**: Snowflake IDs as `int64` for uniqueness across distributed systems
- **API Format**: IDs formatted as strings in JSON responses using `utils.FormatID()`
- **Parsing**: Convert string IDs using `utils.ParseID()` with validation
- **Database**: Store as `bigint` in PostgreSQL for performance

### Time Handling
- **Database**: Use `pgtype.Timestamptz` for timezone-aware timestamps
- **API**: ISO 8601 format in JSON responses
- **Business Logic**: Time slot validation using `common.ParseTimeSlot()`

### Response Structure
```go
// Success responses
c.JSON(http.StatusOK, common.SuccessResponse(data))
c.JSON(http.StatusCreated, common.SuccessResponse(data))

// Error responses (handled by error middleware)
errorCodes.AbortWithError(c, errorCodes.ErrorCode, details)
errorCodes.RespondWithServiceError(c, serviceError)
```

### Testing Patterns
- **Handler Tests**: Test HTTP layer with mock services
- **Service Tests**: Test business logic with mock repositories  
- **Repository Tests**: Test data access with real database (when needed)
- **Mock Interfaces**: Define interfaces for all external dependencies
- **SQLX Mocks**: Use `testutils/mocks/repository.go` for dynamic query mocks

### Memories and Notes
- When generate API, do not forget to add route in @cmd/server/main.go 