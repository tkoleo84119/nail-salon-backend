# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Workflow

### Task Planning and Execution
**CRITICAL - MANDATORY WORKFLOW**: Before starting ANY task, you MUST follow this workflow exactly:

1. **Task Analysis**: Read and understand the complete request
2. **Create Todo List**: Use TodoWrite tool to break down the task into specific, actionable items
3. **Present Plan**: Show the user the planned steps and wait for confirmation before proceeding
4. **Execute Step by Step**: Work through each todo item systematically
5. **Commit Granularly**: Make small, focused commits for each logical change (avoid large commits)
6. **Update Documentation**: After task completion, check if CLAUDE.md needs updates

**⚠️ IMPORTANT**: You MUST NOT start coding or making changes without first creating a todo list and getting user approval. This is non-negotiable and prevents scope creep and ensures proper task planning.

### Commit Guidelines
- **Small, Focused Commits**: Each commit should represent one logical change
- **Clear Messages**: Use conventional commit format without AI-generated indicators
- **Incremental Progress**: Prefer multiple small commits over one large commit
- **Test After Each Change**: Ensure tests pass before committing

### Documentation Maintenance
- **Always Check**: After completing any task, evaluate if CLAUDE.md needs updates
- **Keep Current**: Update architecture decisions, patterns, and best practices
- **Document New Patterns**: Add any new development patterns or conventions discovered

## Common Development Commands

### Running the Application
```bash
make run                    # Start the server (default port from config)
go run cmd/server/main.go   # Alternative way to start server
```

### Testing
```bash
make test                   # Run all tests
go test ./...               # Alternative way to run tests
```

### Database Operations
```bash
make migrate-up             # Run all pending migrations
make migrate-down NUMBER=1  # Rollback migrations (specify number)
make seed-test              # Load test data into database
sqlc generate               # Generate Go code from SQL queries
```

Database migrations are located in `migration/` directory and follow numbered naming convention (001, 002, etc.).

### Environment Setup
- Copy `.env.example` to `.env` and configure database credentials
- **Required environment variables**:
  - `JWT_SECRET` (for JWT token signing)
  - Either `DB_DSN` (connection string) OR individual DB settings: `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_PORT`, `DB_HOST`
- **Optional environment variables**:
  - `PORT` (server port, defaults to 3000)
  - `JWT_EXPIRY_HOURS` (JWT expiry, defaults to 1 hour)
  - `SNOWFLAKE_NODE_ID` (for distributed ID generation, defaults to 1)
  - `DB_SSLMODE` (SSL mode, defaults to disable)
- **Advanced Database Settings**:
  - `DB_MAX_OPEN_CONNS` (max open connections, defaults to 25)
  - `DB_CONN_MAX_LIFE` (connection max lifetime, defaults to 30m)
  - `DB_CONN_MAX_LIFE_JITTER` (connection lifetime jitter, defaults to 5m)
  - `DB_CONN_IDLE_TIME` (connection idle time, defaults to 2m)

## Architecture Overview

### Project Structure
This is a Go-based nail salon management backend using **Clean Architecture** principles with **Domain-Driven Design** elements:

```
nail-salon-backend/
├── cmd/server/                 # Application entry point
├── internal/
│   ├── config/                # Configuration management with environment loading
│   ├── errors/                # Centralized error management system
│   │   ├── codes.go          # Error code constants
│   │   ├── error_manager.go  # YAML-based error manager
│   │   ├── errors.yaml       # Error definitions with Chinese messages
│   │   ├── response_helper.go # Error response utilities
│   │   └── service_error.go  # Service error handling
│   ├── handler/               # HTTP request handlers (Presentation Layer)
│   │   └── staff/            # Staff module handlers
│   ├── infra/db/             # Database infrastructure layer
│   ├── middleware/           # HTTP middleware
│   ├── model/                # Domain models organized by business module
│   │   ├── common/          # Shared API response structures
│   │   └── staff/           # Staff domain models and constants
│   ├── repository/sqlc/      # Data Access Layer with sqlc
│   │   ├── dbgen/           # Generated type-safe queries (centralized)
│   │   └── *.sql            # SQL query files organized by domain
│   ├── service/             # Business Logic Layer
│   │   └── staff/          # Staff business logic services
│   └── utils/               # Shared utilities (JWT, passwords, IDs, etc.)
├── migration/               # Sequential database migration files
├── scripts/seed/           # Database seeding utilities
└── docs/                  # API documentation and database schema
```

### Technology Stack
- **Web Framework**: Gin (github.com/gin-gonic/gin) for high-performance HTTP handling
- **Database**: PostgreSQL with **unified driver strategy**:
  - **pgx/v5**: Unified PostgreSQL driver for all database operations
  - **pgx/v5/stdlib**: Standard library interface for sqlx compatibility
  - **sqlx**: For complex dynamic queries when sqlc limitations are reached
- **Query Generation**: sqlc for type-safe SQL operations with manual ID insertion
- **ID Generation**: Snowflake (github.com/bwmarrin/snowflake) for distributed unique IDs
- **Authentication**: JWT (github.com/golang-jwt/jwt/v5) with refresh token rotation
- **Password Security**: golang.org/x/crypto for bcrypt hashing
- **Environment**: godotenv for local development environment management
- **Testing**: testify with comprehensive mock support

### Clean Architecture Implementation

**Layer Dependencies (Dependency Inversion):**
```
Presentation (handlers/) → Business Logic (service/) → Data Access (repository/) → Infrastructure (infra/)
```

**Key Patterns:**
- **Dependency Injection**: Services accept `dbgen.Querier` interface for testability
- **Interface Segregation**: Clear contracts between layers
- **Modular Organization**: Business domains organized as modules (`staff/`, future `customer/`, `booking/`)
- **Single Responsibility**: Each layer has distinct responsibilities

### Database Design Philosophy
The system manages a comprehensive nail salon business with these design principles:

**Core Design Decisions:**
- **Distributed-Ready IDs**: Snowflake-generated bigint primary keys for scalability
- **Multi-Tenancy**: Store-based access control for multi-location businesses
- **Audit Trails**: Consistent `created_at`/`updated_at` timestamps
- **Domain Modeling**: Rich business model covering all salon operations

**Main Business Domains:**
1. **Store & Staff Management** - Multi-store support with role-based access
2. **Customer Management** - Customer profiles with authentication preferences
3. **Stylist & Scheduling** - Stylist profiles and time slot management
4. **Booking & Services** - Appointment system with service catalog
5. **Inventory & Products** - Product management with stock tracking
6. **Financial Management** - Expense tracking, checkouts, and account transactions

### Authentication & Authorization Architecture
- **JWT-based Authentication** for staff users with configurable expiry
- **Refresh Token System** with device tracking (user agent, IP address)
- **Role-based Access Control** with four defined roles:
  ```go
  RoleSuperAdmin = "SUPER_ADMIN"  // Access to all stores
  RoleAdmin      = "ADMIN"        // Store-specific admin rights
  RoleManager    = "MANAGER"      // Store management capabilities
  RoleStylist    = "STYLIST"      // Service provider access
  ```
- **Store-level Permissions**: SUPER_ADMIN gets all stores, others get explicit assignments
- **Security Features**: Password hashing, token rotation, device tracking

## Development Guidelines

### Code Organization Principles

**Modular Architecture by Business Domain:**
- Organize code by business modules (`staff/`, `customer/`, `booking/`)
- Each module contains: `handler/`, `service/`, `model/` subdirectories
- Shared utilities in `/utils/` for cross-module functionality

**Dependency Management:**
- Use interfaces for all service dependencies
- Pass `dbgen.Querier` interface to services for standard database operations
- Use `sqlx.DB` for dynamic query operations requiring conditional SQL building
- Store configuration in `internal/config/config.go` with centralized env var loading
- Follow dependency injection pattern with interface compliance verification

### Database Operations

**Migrations:**
- Use numbered sequential migration files (001, 002, etc.) in `migration/` directory
- Each migration has separate `.up.sql` and `.down.sql` files
- Follow domain organization for related schema changes

**Query Generation:**
- **Primary**: Use sqlc for type-safe database operations
  - SQL files organized by module in `internal/repository/sqlc/[module]/`
  - Generated code centralized in `internal/repository/sqlc/dbgen/`
  - Run `sqlc generate` after adding/modifying SQL queries
- **Dynamic Operations**: Use sqlx for complex dynamic operations when sqlc is insufficient
  - Repository structure: `internal/repository/sqlx/[module].go`
  - Example: Update operations with optional fields, conditional WHERE clauses
  - Pattern: Build dynamic SQL with conditional parts based on provided parameters

**ID Generation Strategy:**
- **Manual ID insertion**: All database inserts use Snowflake-generated IDs
- **Utility Functions**: Use `utils.GenerateID()` for ID generation
- **Benefits**: Distributed system ready, time-sortable, no database round-trips

### Testing Strategy

**Comprehensive Testing Approach:**
- **Unit Tests**: Service layer business logic with mock dependencies
- **Integration Tests**: Handler testing with mocked services
- **Mock Verification**: Use testify/mock with interface compliance checks
- **Error Scenarios**: Test both success and failure paths
- **Test Organization**: Mirror production code structure in test files

**Mock Patterns:**
```go
// Service interface compliance
var _ staffService.LoginServiceInterface = (*MockLoginService)(nil)

// Handler testing with mocked services
mockService := new(MockLoginService)
handler := NewLoginHandler(mockService)
```

### Security Practices

**Authentication Security:**
- Passwords hashed with bcrypt (cost factor handled by library)
- JWT secrets read from configuration, never hardcoded
- Refresh tokens include device tracking for security auditing
- Token expiry enforced both in JWT claims and database records

**Authorization Patterns:**
- Role constants prevent hardcoding strings throughout codebase
- Store access validated at service layer
- Input validation with Gin's ShouldBindJSON
- Error message sanitization to prevent information leakage

### Configuration Management

**Environment-based Configuration:**
```go
type Config struct {
    DB     DBConfig      // Database connection and pooling
    JWT    JWTConfig     // Authentication configuration
    Server ServerConfig  // Server settings and Snowflake node ID
}
```

**Configuration Loading:**
- Development: `.env` file with godotenv
- Production: Environment variables
- Centralized in `config.Load()` with validation
- Required vs optional variables clearly defined

### Error Handling Philosophy

The application uses a **centralized error management system** with YAML-based error definitions:

**Error Management Architecture:**
- **Error Codes**: Constants defined in `internal/errors/codes.go`
- **Error Definitions**: YAML configuration in `internal/errors/errors.yaml` with Chinese messages
- **Error Manager**: Singleton pattern for loading and managing error definitions
- **Response Helpers**: Standardized error response generation

**Error Handling Patterns:**
- Service errors wrapped with specific error codes
- Database connection validation with proper cleanup
- Graceful degradation with meaningful error messages
- Consistent error responses across API endpoints using `errorCodes.RespondWithError()`
- **Development-friendly debugging**: In development mode, API responses include `dev_details` field with technical error information
- **Production security**: In production mode (`GIN_MODE=release`), only user-friendly messages are returned

## Error Management System

### YAML-Based Error Definitions
The application uses a centralized error management system with YAML-based error definitions:

**Key Components:**
- **Error Codes**: Constants in `internal/errors/codes.go` for easy reference
- **Error Definitions**: YAML file `internal/errors/errors.yaml` with Chinese messages and HTTP status codes
- **Error Manager**: Singleton pattern for loading and managing error definitions
- **Response Helpers**: Standardized utilities for error responses

**Error Categories:**
- **AUTH**: Authentication and authorization errors
- **USER**: User management and validation errors
- **VAL**: Input validation and format errors
- **SYS**: System and infrastructure errors

**Usage Pattern:**
```go
// Load error definitions (in main.go)
errorManager := errorCodes.GetManager()
if err := errorManager.LoadFromFile("internal/errors/errors.yaml"); err != nil {
    log.Fatalf("Failed to load error definitions: %v", err)
}

// Use in handlers
errorCodes.RespondWithError(c, errorCodes.AuthInvalidCredentials, nil)
errorCodes.RespondWithServiceError(c, err) // Automatically includes dev_details in development

// Custom error with details
errorCodes.RespondWithErrorDetails(c, errorCodes.SysDatabaseError, nil, "connection failed: timeout")

// Direct error manager usage (unified method)
errorManager.GetErrorResponse(code, fieldErrors)                    // Without details
errorManager.GetErrorResponse(code, fieldErrors, "debug info")     // With details
```

**Development vs Production Error Responses:**

Development Mode (GIN_MODE != "release"):
```json
{
  "message": "資料庫連接錯誤",
  "dev_details": "failed to get store access: sql: connection timeout"
}
```

Production Mode (GIN_MODE="release"):
```json
{
  "message": "資料庫連接錯誤"
}
```

## API Design Standards

### Standardized Response Format
All APIs (except `/health`) use a consistent response structure:

**Success Response (2xx):**
```json
{
  "data": {
    // Actual response data
  }
}
```

**Error Response (4xx/5xx):**
```json
{
  "message": "錯誤描述訊息",
  "errors": {
    "field_name": "具體錯誤說明"
  }
}
```

### JSON Response Format Standards
**Field Naming Convention:**
- All JSON response field names use **camelCase** (e.g., `userId`, `storeId`, `storeList`)
- Database fields remain snake_case internally for consistency with PostgreSQL conventions

**ID Field Format:**
- All ID fields in API responses are returned as **strings** (e.g., `"123456789"`)
- Frontend requests should send ID fields as strings
- Backend internally converts between string (API) and int64 (database) using `utils.ParseID()` and `utils.FormatID()`

**Field Conversion Utilities:**
- `utils.ParseID(string)` - Convert string ID to int64 for database operations
- `utils.FormatID(int64)` - Convert int64 ID to string for API responses
- `utils.ParseIDSlice([]string)` - Convert string ID slice to int64 slice
- `utils.FormatIDSlice([]int64)` - Convert int64 ID slice to string slice

### Error Message Localization
- All error messages are in Traditional Chinese
- Field validation errors include specific field names in Chinese
- Use `internal/utils/validation.go` for extracting and translating validation errors
- Common field translations: Username→帳號, Password→密碼, StoreID→門市ID

### Response Structure Guidelines
- **Success responses**: Only include `data` field, no `message`
- **Validation errors**: Use `輸入驗證失敗` as message with detailed field errors
- **Authentication errors**: Use `認證失敗` with credentials error details
- **System errors**: Use `系統錯誤` with server error information
- **Request format errors**: Use `請求錯誤` with request parsing errors

### Implementation Pattern
```go
// Success
c.JSON(http.StatusOK, common.SuccessResponse(responseData))

// Service Error
errorCodes.RespondWithServiceError(c, err)

// Validation Error
errorCodes.RespondWithError(c, errorCodes.ValInputValidationFailed, validationErrors)

// Custom Error
errorCodes.RespondWithError(c, errorCodes.AuthTokenMissing, nil)

// Middleware Error (with abort)
errorCodes.AbortWithError(c, errorCodes.AuthPermissionDenied, nil)
```

## API Documentation Standards
API documentation maintained in `docs/api/` directory using structured format:
- Request/response examples with standardized format
- Error response documentation in Chinese
- Authentication requirements
- Business logic explanations

## Commit Message Convention
This project follows Conventional Commits specification:
```
<type>: <description>

Types: feat, fix, refactor, perf, style, test, docs, build, ops, chore
- Use English for descriptions
- Keep descriptions under 100 characters
- Use imperative mood (e.g., "add feature" not "added feature")
- Include context about architectural decisions when relevant
- Do not include AI-generated content indicators
```

## Important Development Reminders

**DO:**
- **ALWAYS start tasks by creating a todo list and getting user confirmation** - This is the single most important rule
- Follow the modular architecture pattern established by the staff module
- Use role constants instead of hardcoded strings
- Generate Snowflake IDs for all database insertions
- Write comprehensive tests for both success and failure scenarios
- Use sqlc for type-safe database operations as primary approach
- Use sqlx for dynamic operations requiring conditional SQL (e.g., optional field updates)
- Centralize configuration in config layer
- Follow Clean Architecture dependency flow
- Use standardized API response format with `common.SuccessResponse()` and `errorCodes.RespondWithError()`
- Provide Chinese error messages for all user-facing errors
- Extract validation errors using `utils.ExtractValidationErrors()`
- Make small, focused commits with clear messages
- Update CLAUDE.md after completing tasks that introduce new patterns
- Load error definitions from YAML file in main.go initialization
- Use the centralized error management system for all error responses
- Prefer `RespondWithServiceError()` for service layer errors (includes development details automatically)
- Use `RespondWithErrorDetails()` when you want to add specific debug information

**DON'T:**
- **Start coding without presenting a plan to the user first** - This will result in incomplete/incorrect implementations
- Put business logic in handlers (belongs in service layer)
- Read environment variables directly with `os.Getenv` outside config layer
- Use auto-incrementing database IDs (use Snowflake IDs)
- Hardcode role strings (use constants from `staff.Role*`)
- Skip testing error scenarios
- Mix dynamic queries with sqlc unless necessary (use sqlx for dynamic operations)
- Return inconsistent API response formats (always use camelCase and string IDs)
- Use English error messages for user-facing errors
- Mix snake_case and camelCase in JSON responses (always use camelCase)
- Return ID fields as numbers in API responses (always use strings)
- Hardcode validation error messages in handlers
- Make large commits that change too many things at once
- Use different PostgreSQL drivers (stick to pgx/v5)
- Modify error messages directly in code (use `internal/errors/errors.yaml`)
- Skip loading error definitions in main.go (required for error system to work)

**ARCHITECTURE NOTES:**
- The system is designed for multi-store nail salon businesses
- Authentication supports both staff users and customers (different flows)
- Store-level access control enables franchise/multi-location scenarios
- Snowflake IDs prepare the system for distributed deployment
- The modular approach supports incremental feature development