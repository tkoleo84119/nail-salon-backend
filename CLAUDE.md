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

## Authentication System

### Dual Authentication Model
The system supports two distinct authentication contexts:

**Customer Authentication (Public API)**
- LINE OAuth integration for customer login
- JWT tokens with CustomerContext
- Routes: `/api/auth/line/login`, `/api/auth/line/register`
- Middleware: `CustomerJWTAuth`
- Context extraction: `GetCustomerFromContext()`

**Staff Authentication (Admin API)**
- Traditional username/password authentication
- JWT tokens with StaffContext including roles and store access
- Routes: `/api/admin/auth/login`, `/api/admin/auth/token/refresh`
- Middleware: `JWTAuth` with role-based authorization
- Hierarchical roles: `SUPER_ADMIN` > `ADMIN` > `MANAGER` > `STYLIST`

### Role-Based Access Control (only for admin API)
- `RequireAdminRoles()` - SUPER_ADMIN, ADMIN only
- `RequireManagerOrAbove()` - MANAGER, ADMIN, SUPER_ADMIN
- `RequireAnyStaffRole()` - Any authenticated staff member
- `RequireRoles(roles...)` - Specific role combinations

## Database Architecture

### Triple Connection Strategy
```go
type Database struct {
    Std     *sql.DB        // Standard library compatibility
    Sqlx    *sqlx.DB       // SQLX for complex operations
    PgxPool *pgxpool.Pool  // PGX for high performance
}
```

### SQLC Integration
- Type-safe SQL queries generated from `.sql` files in `/internal/repository/sqlc/`
- Configuration in `sqlc.yaml` with JSON tags and interfaces enabled
- Generated code provides `Querier` interface for unified database operations
- PostgreSQL-specific with pgtype for Go type conversion

### Repository Pattern
- **SQLC repositories**: Generated CRUD operations with type safety
- **SQLX repositories**: Complex dynamic queries and partial updates
- **Interface-based design**: Enables testing and maintainability
- **Dependency injection**: Services receive repository interfaces

## API Design Patterns

### Route Structure
**Public/Customer Routes** (`/api/`)
```
/api/auth/line/login              # LINE OAuth authentication
/api/customers/me                 # Customer profile management
/api/bookings                     # Customer booking operations
/api/stores                       # Store browsing and services
/api/stores/:id/stylists          # Store stylist information
/api/stores/:id/stylists/:id/schedules  # Schedule browsing
/api/schedules/:id/time-slots     # Available time slot viewing
```

**Admin Routes** (`/api/admin/`)
```
/api/admin/auth/login             # Staff authentication
/api/admin/staff                  # Staff management
/api/admin/stores                 # Store management
/api/admin/services               # Service management
/api/admin/schedules/bulk         # Bulk schedule operations
/api/admin/time-slot-templates    # Schedule template management
```

### Service Layer Organization
Services are organized by access level and domain:

```go
Services {
    Public PublicServices    # Customer-facing business logic
    Admin  AdminServices     # Staff-facing business logic
}
```

Each service follows the pattern:
- **Interface definition** for testability
- **Constructor function** with dependency injection
- **Context-aware methods** for all operations
- **ServiceError returns** for consistent error handling

## Error Handling

### Error Code System
Hierarchical error codes for different domains:
- `AUTH_*` - Authentication and authorization errors
- `BOOKING_*` - Booking-related business logic errors
- `CUSTOMER_*` - Customer management errors
- `SCHEDULE_*` - Scheduling system errors
- `SERVICE_*` - Service management errors
- `STORE_*` - Store management errors
- `VAL_*` - Validation errors
- `SYS_*` - System-level errors

### Error Flow
```go
// Service layer generates business logic errors
ServiceError{Code: "CUSTOMER_NOT_FOUND", Message: "Customer not found"}

// Handler layer maps to HTTP responses
errorCodes.RespondWithServiceError(c, err)

// Client receives structured error response
{"error": "CUSTOMER_NOT_FOUND", "message": "Customer not found"}
```

## Development Guidelines

### Code Style
- Follow Go conventions and gofmt formatting
- Use meaningful package aliases for clarity (e.g., `adminAuthModel`, `customerModel`)
- Interface-first design for services and repositories
- Comprehensive error handling with specific error codes
- Context propagation throughout the application
- Dependency injection through container pattern

### Database Guidelines
- Use SQLC for standard CRUD operations and complex queries
- Use SQLX for dynamic queries and partial updates
- Always use parameterized queries to prevent SQL injection
- Implement proper transaction handling for multi-step operations
- Use pgtype for PostgreSQL-specific data types (Time, Date, etc.)

### Security Best Practices
- JWT token validation on all protected routes
- Role-based authorization for admin operations
- Customer blacklist validation before booking operations
- Store access control based on staff assignments
- Input validation using struct tags and custom validators
- Secure password hashing with bcrypt

### Testing Strategy
- Unit tests for all utility functions
- Service layer testing with interface mocking
- Table-driven tests for comprehensive scenario coverage
- Integration tests for database operations
- Use testify framework for assertions and test structure

## Key Business Logic

### Booking System
- Customers can create, view, update, and cancel bookings
- Time slot availability checking with booking status validation
- Store and stylist assignment with access control
- Blacklist validation prevents bookings from banned customers

### Schedule Management
- Bulk schedule creation and deletion by staff
- Time slot templates for efficient schedule setup
- Individual time slot management with duration tracking
- Available time slot filtering for customer browsing

### Store Management
- Multi-store support with store-specific staff access
- Service offerings per store with pricing
- Stylist assignment and availability per store
- Store access control based on staff roles

### Customer Management
- LINE social login integration
- Customer profile management
- Booking history and status tracking
- Blacklist management for access control

## Environment Configuration

Required environment variables:
- `DB_URL` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret
- `LINE_CHANNEL_ID` - LINE OAuth channel ID
- `LINE_CHANNEL_SECRET` - LINE OAuth channel secret
- `PORT` - Server port (default: 8080)

## Development Workflow

When implementing new features:
- Define models in appropriate `/internal/model/` subdirectory
- Create SQLC queries in `/internal/repository/sqlc/` if needed
- Implement service layer with interface and concrete implementation
- Create handler with proper error handling and validation
- Add routes to appropriate router setup function
- Update dependency injection in container files
- Write comprehensive tests for new functionality
- Update API documentation in `/docs/` directory

When modifying existing features:
- Understand the layered architecture before making changes
- Follow established patterns for consistency
- Ensure interface compatibility when updating services
- Run full test suite to verify no regressions
- Update documentation if API contracts change

## Common Patterns

### Time Handling
- Use `time.Time` for Go operations
- Use `pgtype.Time` for database storage
- Convert between formats as needed for business logic
- Handle timezone considerations for scheduling

### ID Generation
- Use Snowflake algorithm for unique IDs across all entities
- Initialize node ID based on environment or configuration
- IDs are int64 but often represented as strings in JSON

### Validation
- Use struct tags for basic validation (`binding:"required"`)
- Implement custom validators for business rules (phone numbers, dates)
- Extract validation errors for user-friendly responses
- Validate business logic constraints in service layer

### Pagination
- Implement cursor-based pagination for better performance
- Use limit/offset for simple cases
- Provide total counts when feasible
- Sort by creation time or relevant business criteria