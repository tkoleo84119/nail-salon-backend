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
├── app/             # Application layer (dependency injection, routing)
│   ├── container.go         # Main dependency injection container
│   ├── container_admin.go   # Admin-specific services/handlers
│   ├── container_public.go  # Public/customer services/handlers
│   └── router.go           # Modular route setup by domain
├── config/          # Environment configuration
├── errors/          # Centralized error management with YAML definitions
├── handler/         # HTTP handlers organized by domain
│   ├── admin/      # Admin-specific handlers (staff, stores, etc.)
│   ├── auth/       # Authentication handlers
│   ├── booking/    # Booking management handlers
│   └── customer/   # Customer-facing handlers
├── infra/db/        # Database connection setup
├── middleware/      # JWT auth and role authorization
├── model/           # Request/response models organized by domain
│   ├── admin/      # Admin-specific models
│   ├── auth/       # Authentication models
│   ├── booking/    # Booking models
│   ├── customer/   # Customer models
│   └── common/     # Shared models and utilities
├── repository/      # Data access layer
│   ├── sqlc/       # SQLC generated queries
│   └── sqlx/       # SQLX manual queries for dynamic operations
├── service/         # Business logic layer organized by domain
│   ├── admin/      # Admin-specific services
│   ├── auth/       # Authentication services
│   ├── booking/    # Booking management services
│   └── customer/   # Customer services
├── testutils/       # Test utilities and mocks
└── utils/           # Utility functions (JWT, passwords, validation, type conversion)
```

### Key Patterns
- **Layered Architecture**: Handler → Service → Repository with clear domain separation
- **Domain-Separated Architecture**: Admin and public APIs with separate containers
- **Dependency Injection**: Multi-container system with domain-specific organization
- **Route Organization**: Modular route setup by business domain with role-based middleware
- **Error Handling**: Centralized error codes with standardized response helpers
- **Authentication**: Dual JWT system (staff and customer) with context extraction
- **Database**: Strategic SQLC/SQLX usage based on operation complexity
- **Type Safety**: Comprehensive type conversion utilities for PostgreSQL integration
- **Testing**: Domain-separated testing with comprehensive mocks

### Database Architecture
- **PostgreSQL**: Primary database with connection pooling via `pgxpool.Pool`
- **pgx/v5**: Primary driver for transactions and connection management
- **SQLC**: Type-safe SQL query generation for standard CRUD operations
- **SQLX**: Dynamic SQL queries for complex updates with optional fields, or complex queries for search
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
Comprehensive nail salon management system with:
- **Staff Management**: Multi-role system (SUPER_ADMIN, ADMIN, MANAGER, STYLIST)
- **Store Management**: Multi-location support with store-specific permissions
- **Customer Management**: Customer accounts and authentication
- **Booking System**: Appointment scheduling with time slot management
- **Service Management**: Service catalog and scheduling
- **Authentication**: Dual authentication system for staff and customers
- **Role-Based Authorization**: Granular permissions based on staff roles

### Commit Convention
Follow Conventional Commits format: `<type>: <description>`
- Types: feat, fix, refactor, perf, style, test, docs, build, ops, chore
- Use imperative mood, present tense, no capitalization, no period
- English language for descriptions

### Development Notes
- If needs new SQL query, look up current SQL query first (in `internal/repository/sqlc/` and `internal/repository/sqlx/`), not create duplicate SQL
- If sql is update optional columns(dynamic), use sqlx not sqlc
- When adding a new SQL for batch insert operations, use sqlc's :copyfrom
- Handler & service & model need to use modules struct by business domain (separate by folder)
- Always use type conversion utilities from `internal/utils/type_convert.go` instead of manual pgtype construction
- For nullable database fields, use appropriate converter functions (e.g., `StringPtrToPgText`, `Int64ToPgInt8`)

### Type Conversion Utilities

The `internal/utils/type_convert.go` provides comprehensive type conversion utilities for PostgreSQL types:

#### String Conversions
- `StringPtrToPgText(s *string, emptyAsNull bool)` - Unified function for optional string fields
  - `emptyAsNull=false`: Empty strings remain as valid empty strings
  - `emptyAsNull=true`: Empty strings are treated as NULL
- `PgTextToString(t pgtype.Text)` - Convert pgtype.Text to string (handles NULL as empty string)

#### Numeric Conversions
- `Float64ToPgNumeric(f float64)` - Convert float64 to pgtype.Numeric (with error handling)
- `Int64ToPgNumeric(i int64)` - Convert int64 to pgtype.Numeric (with error handling)
- `PgNumericToFloat64(n pgtype.Numeric)` - Convert pgtype.Numeric to float64 (handles NULL as 0)

#### Boolean Conversions
- `BoolPtrToPgBool(b *bool)` - Convert bool pointer to pgtype.Bool (nil becomes NULL)

#### Time Conversions
- `TimeToPgTimestamptz(t time.Time)` - Convert time to pgtype.Timestamptz
- `TimeToPgTime(t time.Time)` - Convert time to pgtype.Time
- `TimeToPgDate(t time.Time)` - Convert time to pgtype.Date
- `DateStringToTime(s string)` - Parse YYYY-MM-DD format to time.Time
- `TimeStringToTime(s string)` - Parse HH:MM format to time.Time
- `PgTimeToTimeString(t pgtype.Time)` - Convert pgtype.Time to HH:MM string
- `PgDateToDateString(d pgtype.Date)` - Convert pgtype.Date to YYYY-MM-DD string

#### ID Conversions (for Foreign Keys)
- `Int64ToPgInt8(id int64)` - Convert int64 ID to pgtype.Int8 (for required foreign keys)
- `Int64PtrToPgInt8(id *int64)` - Convert int64 pointer to pgtype.Int8 (for optional foreign keys)
- `Int32ToPgInt4(value int32)` - Convert int32 to pgtype.Int4
- `Int32PtrToPgInt4(value *int32)` - Convert int32 pointer to pgtype.Int4
- `ParseIDToPgInt8(idStr string)` - Parse string ID to pgtype.Int8 with validation
- `ParseIDPtrToPgInt8(idStr *string)` - Parse string ID pointer to pgtype.Int8 with validation
- `PgInt8ToIDString(id pgtype.Int8)` - Convert pgtype.Int8 to ID string (handles NULL)
- `PgInt8ToIDStringPtr(id pgtype.Int8)` - Convert pgtype.Int8 to ID string pointer (NULL returns nil)

#### Usage Guidelines
- **Never manually construct pgtype structures** - Always use conversion utilities
- **Check function documentation** - Each function has comprehensive examples and usage patterns
- **Handle errors properly** - Functions that can fail return errors (numeric conversions, time parsing)
- **Choose appropriate null handling** - Use `emptyAsNull` parameter for string fields based on business logic

### Validation Patterns

#### Handler Layer Validation (Fixed Order)
1. **Input JSON Validation** - Always first, using `c.ShouldBindJSON(&req)`
2. **Path Parameter Validation** - Validate required URL parameters
3. **Business Logic Validation** - Check if update requests have fields using `req.HasUpdates()` or specialized methods
4. **Authentication Context Validation** - Use `middleware.GetStaffFromContext(c)` or `middleware.GetCustomerFromContext(c)`
5. **ID Parsing Validation** - Convert string IDs using `utils.ParseID()` and `utils.ParseIDSlice()`
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
- **Error Handling**: Use `errorCodes.AbortWithError()` in handlers, `errorCodes.RespondWithServiceError()` for service errors
- **Standardized Helpers**: Use `errorCodes.RespondWithEmptyFieldError(c)` for update validation failures
- **Context Extraction**: Use `middleware.GetStaffFromContext(c)` for staff, `middleware.GetCustomerFromContext(c)` for customers
- **ID Parsing**: Use `utils.ParseID()` and `utils.ParseIDSlice()` for string to int64 conversion with validation
- **Time Validation**: Use `common.ParseTimeSlot()` for time format validation
- **Update Validation**: Implement `HasUpdates()` and specialized validation methods on request models

### Error Handling Patterns
- **Centralized Error Management**: Error codes defined in `internal/errors/errors.yaml`
- **Service Errors**: `errorCodes.NewServiceError()` and `errorCodes.NewServiceErrorWithCode()`
- **Handler Errors**: `errorCodes.AbortWithError()` and `errorCodes.RespondWithServiceError()`
- **Standardized Helpers**: `errorCodes.RespondWithEmptyFieldError(c)` for consistent empty field responses
- **Error Categories**: AUTH, BOOKING, CUSTOMER, SCHEDULE, SERVICE, STORE, USER, VAL, SYS
- **Permission Errors**: Use `AUTH.AUTH_PERMISSION_DENIED` for all permission-related issues
- **Localized Messages**: Error messages in Chinese for user-facing responses

## CRUD Operation Patterns

### CREATE Operations

#### Handler Pattern (Create)
```go
func (h *Handler) CreateEntity(c *gin.Context) {
    // Input JSON validation
    var req EntityCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErrors := utils.ExtractValidationErrors(err)
        errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
        return
    }

    // Authentication context validation
    staffContext, exists := middleware.GetStaffFromContext(c)
    if !exists {
        errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
        return
    }

    // Service layer call
    response, err := h.service.CreateEntity(c.Request.Context(), req, *staffContext)
    if err != nil {
        errorCodes.RespondWithServiceError(c, err)
        return
    }

    // Success response (201 Created)
    c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
```

#### Service Pattern (Create)
```go
func (s *Service) CreateEntity(ctx context.Context, req EntityRequest, staffContext StaffContext) (*EntityResponse, error) {
    // Input validation & ID parsing
    staffUserID, err := utils.ParseID(staffContext.UserID)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "invalid staff user ID", err)
    }

    // Business logic validation (permissions, role)
    if err := s.validatePermissions(staffContext.Role, req); err != nil {
        return nil, err
    }

    // Data integrity validation (existence, uniqueness)
    exists, err := s.queries.CheckEntityExists(ctx, req.Name)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "check failed", err)
    }
    if exists {
        return nil, errorCodes.NewServiceErrorWithCode(errorCodes.EntityAlreadyExists)
    }

    // Generate entity ID
    entityID := utils.GenerateID()

    // Transaction-based creation
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "transaction failed", err)
    }
    defer tx.Rollback(ctx)

    qtx := dbgen.New(tx)

    // Create entity using type converters
    created, err := qtx.CreateEntity(ctx, dbgen.CreateEntityParams{
        ID:          entityID,
        Name:        req.Name,
        Description: utils.StringPtrToPgText(&req.Description, false), // Empty strings allowed
        Note:        utils.StringPtrToPgText(&req.Note, true),        // Empty as NULL
        IsActive:    utils.BoolPtrToPgBool(&req.IsActive),
        CreatedBy:   utils.Int64ToPgInt8(staffUserID),
    })
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "create failed", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "commit failed", err)
    }

    return &EntityResponse{
        ID:          utils.FormatID(created.ID),
        Name:        created.Name,
        Description: created.Description.String,
        Note:        created.Note.String,
        IsActive:    created.IsActive.Bool,
    }, nil
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
    // Path parameter validation
    targetID := c.Param("id")
    if targetID == "" {
        errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
            "id": "id為必填項目",
        })
        return
    }

    // Input JSON validation
    var req EntityUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErrors := utils.ExtractValidationErrors(err)
        errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
        return
    }

    // Business logic validation - HasUpdates check (use standardized helper)
    if !req.HasUpdates() {
        errorCodes.RespondWithEmptyFieldError(c)
        return
    }

    // Authentication context validation (support both staff and customer)
    staffContext, exists := middleware.GetStaffFromContext(c)
    if !exists {
        errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
        return
    }

    // Service call
    response, err := h.service.UpdateEntity(c.Request.Context(), targetID, req, *staffContext)
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

    // Dynamic field updates using type converters
    if req.Field1 != nil {
        setParts = append(setParts, "field1 = :field1")
        args["field1"] = utils.StringPtrToPgText(req.Field1, true) // Empty as NULL
    }

    if req.Field2 != nil {
        setParts = append(setParts, "field2 = :field2")
        args["field2"] = utils.BoolPtrToPgBool(req.Field2)
    }

    if req.UpdaterID != nil {
        setParts = append(setParts, "updater_id = :updater_id")
        updaterID, err := utils.ParseID(*req.UpdaterID)
        if err != nil {
            return nil, fmt.Errorf("invalid updater ID: %w", err)
        }
        args["updater_id"] = utils.Int64ToPgInt8(updaterID)
    }

    query := fmt.Sprintf(`
        UPDATE entities SET %s WHERE id = :id
        RETURNING id, field1, field2, updater_id, updated_at
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
    // ID parsing and validation
    id, err := utils.ParseID(entityID)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid ID", err)
    }

    // Entity existence check
    entity, err := s.queries.GetEntityByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errorCodes.NewServiceErrorWithCode(errorCodes.EntityNotFound)
        }
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "get failed", err)
    }

    // Permission validation
    if err := s.validateDeletePermissions(staffContext, entity); err != nil {
        return nil, err
    }

    // Business constraint validation
    hasConstraints, err := s.queries.CheckEntityConstraints(ctx, id)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "constraint check failed", err)
    }
    if hasConstraints {
        return nil, errorCodes.NewServiceErrorWithCode(errorCodes.EntityConstraintViolation)
    }

    // Perform delete (usually soft delete)
    if err := s.queries.DeleteEntity(ctx, id); err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "delete failed", err)
    }

    return buildDeleteResponse(entity), nil
}
```

## Model Validation Patterns

### Request Model Structure
```go
// Create requests use direct types for required fields
type CreateEntityRequest struct {
    Name     string   `json:"name" binding:"required,min=1,max=100"`
    Email    string   `json:"email" binding:"required,email"`
    Role     string   `json:"role" binding:"required,oneof=ADMIN MANAGER STYLIST"`
    StoreIDs []string `json:"storeIds" binding:"required,min=1,max=100"`
}

// Update requests use pointer types for optional fields
type UpdateEntityRequest struct {
    Name     *string  `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
    Email    *string  `json:"email,omitempty" binding:"omitempty,email"`
    IsActive *bool    `json:"isActive,omitempty"`
    StoreIDs []string `json:"storeIds,omitempty" binding:"omitempty,min=1,max=100"`
}

// HasUpdates method required for all update requests
func (r UpdateEntityRequest) HasUpdates() bool {
    return r.Name != nil || r.Email != nil || r.IsActive != nil || len(r.StoreIDs) > 0
}

// Advanced validation methods for complex business logic
func (r UpdateBookingRequest) HasTimeSlotUpdate() bool {
    return r.StoreId != nil || r.StylistId != nil || r.TimeSlotId != nil
}

func (r UpdateBookingRequest) IsTimeSlotUpdateComplete() bool {
    if !r.HasTimeSlotUpdate() {
        return true // No time slot update, so it's complete
    }
    return r.StoreId != nil && r.StylistId != nil && r.TimeSlotId != nil
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
// Staff authentication (admin endpoints)
staffContext, exists := middleware.GetStaffFromContext(c)
if !exists {
    errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
    return
}

// Customer authentication (public endpoints)
customerContext, exists := middleware.GetCustomerFromContext(c)
if !exists {
    errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
    return
}
```

### Role-Based Authorization
- **SUPER_ADMIN**: Full system access, can manage all entities across all stores
- **ADMIN**: Store-level administration, can manage assigned stores and their staff
- **MANAGER**: Store operations, limited administrative functions within assigned stores
- **STYLIST**: Self-service only, can only manage own records and view own schedule

### Role-Based Middleware
- `RequireAdminRoles()` - SUPER_ADMIN, ADMIN only
- `RequireManagerOrAbove()` - SUPER_ADMIN, ADMIN, MANAGER
- `RequireAnyStaffRole()` - All staff roles
- `RequireSuperAdmin()` - SUPER_ADMIN only
- `RequireRoles(role1, role2)` - Custom role combinations

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

### Application Architecture Patterns

#### Enhanced Container Architecture

**Main Container** (`internal/app/container.go`):
- **Central Orchestration**: Coordinates all domain-specific containers
- **Shared Dependencies**: Database, configuration, and common utilities
- **Clean Separation**: Clear boundaries between admin and public domains

**Domain-Specific Containers**:
- **Admin Container** (`container_admin.go`): Staff management, store administration
- **Public Container** (`container_public.go`): Customer-facing services, booking

**Container Structure**:
```go
type Container struct {
    cfg      *config.Config
    database *db.Database

    repositories Repositories
    services     Services
    handlers     Handlers
}

type Services struct {
    Public PublicServices // Customer, booking, etc.
    Admin  AdminServices  // Staff, store, admin functions
}

type Handlers struct {
    Public PublicHandlers
    Admin  AdminHandlers
}
```

#### Route Organization (`internal/app/router.go`)
- **Domain Separation**: Clear separation between admin and public API routes
- **Modular Setup Functions**: Each domain has dedicated route setup functions
- **Consistent Middleware**: JWT authentication and role-based authorization applied systematically
- **Role-Based Protection**: Different middleware for different access levels
- **Separated Concerns**: Route logic completely separated from main application bootstrap

**Route Structure**:
- **Admin Routes**: `/admin/*` with staff JWT and role-based middleware
- **Public Routes**: `/api/*` with customer JWT for protected endpoints
- **Auth Routes**: `/auth/*` for authentication (both staff and customer)
- **Health Routes**: `/health` for monitoring

#### Refactored Main (`cmd/server/main.go`)
- **Minimal Bootstrap**: Only essential initialization (config, database, error manager, validators)
- **Clean Separation**: Business logic moved to app layer, main focuses on startup
- **Maintainable**: Easy to understand and modify application startup sequence

### Development Patterns

#### Adding New APIs
1. **Create Handler**: Add new handler in appropriate domain folder
2. **Create Service**: Add business logic in service layer
3. **Update Container**: Add service and handler to container initialization
4. **Update Router**: Add route in appropriate route setup function
5. **No Main Changes**: Routes automatically available without modifying main.go

#### Container Usage Pattern
```go
// Access dependencies through container
container := app.NewContainer(cfg, database)
handlers := container.GetHandlers()
services := container.GetServices()
repositories := container.GetRepositories()
```

### Latest Development Patterns

#### Container Usage Pattern
```go
// Initialize container with domain separation
container := app.NewContainer(cfg, database)
handlers := container.GetHandlers()
services := container.GetServices()
repositories := container.GetRepositories()

// Access domain-specific handlers
adminHandlers := handlers.Admin
publicHandlers := handlers.Public
```

#### Service Layer Pattern (Enhanced)
```go
func (s *Service) Operation(ctx context.Context, req Request, staffContext common.StaffContext) (*Response, error) {
    // 1. Input validation & ID parsing (with slice support)
    staffUserID, err := utils.ParseID(staffContext.UserID)
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "invalid staff user ID", err)
    }

    storeIDs, err := utils.ParseIDSlice(req.StoreIDs) // For slice parsing
    if err != nil {
        return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store IDs", err)
    }

    // 2. Request completeness validation
    if !req.HasUpdates() {
        return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field required", nil)
    }

    // 3. Business logic validation (roles, complex validation)
    if err := s.validateComplexBusinessLogic(req); err != nil {
        return nil, err
    }

    // 4. Permission & authorization validation
    if err := s.validatePermissions(staffContext, req); err != nil {
        return nil, err
    }

    // 5. Data integrity validation
    // 6. Database operations with appropriate transaction usage
}
```

#### Repository Pattern Guidelines
**SQLC Usage (Preferred for)**:
- Standard CRUD operations with fixed parameters
- Existence checks and simple queries
- Batch operations using `:copyfrom` syntax
- Performance-critical queries

**SQLX Usage (Required for)**:
- Dynamic update operations with optional fields
- Complex WHERE clauses with variable conditions
- Flexible field selection in queries

### Current Best Practices

#### Error Handling
- **Handlers**: Use `errorCodes.RespondWithEmptyFieldError(c)` for consistent empty field responses
- **Services**: Use `errorCodes.NewServiceError()` with proper error categories
- **Database Errors**: Use `errors.Is(err, pgx.ErrNoRows)` for no-data checks

#### Authentication & Authorization
- **Dual System**: Separate staff and customer authentication flows
- **Context Extraction**: Always use middleware functions, never direct context access
- **Role Validation**: Use appropriate middleware for different access levels

#### Type Conversion
- **Always use utilities**: Never manually construct pgtype structures
- **Error Handling**: Properly handle conversion errors, especially for numeric types
- **Null Handling**: Choose appropriate `emptyAsNull` parameter for string fields

#### Model Design
- **Complex Validation**: Implement specialized validation methods for business logic
- **Update Requests**: Use pointer types for all optional fields
- **Validation Tags**: Use comprehensive binding tags for all field constraints

### Development Notes
- When generating APIs, add routes in appropriate domain setup functions in `internal/app/router.go`
- Use domain-separated containers for new services and handlers
- Always implement proper validation order in both handlers and services
- For permission issues, use `AUTH.AUTH_PERMISSION_DENIED` error code
- Implement `HasUpdates()` methods for all update request models
- Use transactions appropriately based on operation complexity
