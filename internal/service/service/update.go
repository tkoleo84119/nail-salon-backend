package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateServiceService struct {
	queries     dbgen.Querier
	serviceRepo sqlx.ServiceRepositoryInterface
}

func NewUpdateServiceService(queries dbgen.Querier, serviceRepo sqlx.ServiceRepositoryInterface) *UpdateServiceService {
	return &UpdateServiceService{
		queries:     queries,
		serviceRepo: serviceRepo,
	}
}

func (s *UpdateServiceService) UpdateService(ctx context.Context, serviceID string, req service.UpdateServiceRequest, updaterRole string) (*service.UpdateServiceResponse, error) {
	// Validate permissions
	if err := s.validatePermissions(updaterRole); err != nil {
		return nil, err
	}

	// Parse service ID
	parsedServiceID, err := utils.ParseID(serviceID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid service ID", err)
	}

	// Validate request has at least one field to update
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Check if service exists
	_, err = s.queries.GetServiceByID(ctx, parsedServiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get service", err)
	}

	// Check name uniqueness if name is being updated
	if req.Name != nil {
		exists, err := s.queries.CheckServiceNameExistsExcluding(ctx, dbgen.CheckServiceNameExistsExcludingParams{
			Name: *req.Name,
			ID:   parsedServiceID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check service name uniqueness", err)
		}
		if exists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceAlreadyExists)
		}
	}

	// Update service using sqlx repository
	updatedService, err := s.serviceRepo.UpdateService(ctx, parsedServiceID, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update service", err)
	}

	return updatedService, nil
}

// validatePermissions checks if the updater has permission to update services
func (s *UpdateServiceService) validatePermissions(updaterRole string) error {
	switch updaterRole {
	case staff.RoleSuperAdmin, staff.RoleAdmin:
		return nil
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}
