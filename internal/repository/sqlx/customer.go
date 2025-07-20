package sqlx

import (
	"github.com/jmoiron/sqlx"
)

// CustomerRepositoryInterface defines the interface for customer repository
type CustomerRepositoryInterface interface {
	// Add any dynamic customer operations here if needed in the future
}

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Currently no dynamic operations needed for customer LINE login
// All operations can be handled by SQLC generated queries