package adminService

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateServiceService struct {
	repo *sqlx.Repositories
}

func NewUpdateServiceService(repo *sqlx.Repositories) *UpdateServiceService {
	return &UpdateServiceService{
		repo: repo,
	}
}

func (s *UpdateServiceService) UpdateService(ctx context.Context, serviceID int64, req adminServiceModel.UpdateServiceRequest, updaterRole string) (*adminServiceModel.UpdateServiceResponse, error) {
	// Validate permissions
	if err := s.validatePermissions(updaterRole); err != nil {
		return nil, err
	}

	// Validate request has at least one field to update
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Check if service exists
	_, err := s.repo.Service.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get service", err)
	}

	// Check name uniqueness if name is being updated
	if req.Name != nil {
		exists, err := s.repo.Service.CheckServiceNameExistsExcluding(ctx, serviceID, *req.Name)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check service name uniqueness", err)
		}
		if exists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceAlreadyExists)
		}
	}

	updatedService, err := s.repo.Service.UpdateService(ctx, serviceID, sqlx.UpdateServiceParams{
		Name:            req.Name,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsAddon:         req.IsAddon,
		IsVisible:       req.IsVisible,
		IsActive:        req.IsActive,
		Note:            req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update service", err)
	}

	response := adminServiceModel.UpdateServiceResponse{
		ID:              utils.FormatID(updatedService.ID),
		Name:            updatedService.Name,
		Price:           int64(utils.PgNumericToFloat64(updatedService.Price)),
		DurationMinutes: updatedService.DurationMinutes,
		IsAddon:         utils.PgBoolToBool(updatedService.IsAddon),
		IsVisible:       utils.PgBoolToBool(updatedService.IsVisible),
		IsActive:        utils.PgBoolToBool(updatedService.IsActive),
		Note:            utils.PgTextToString(updatedService.Note),
		CreatedAt:       utils.PgTimestamptzToTimeString(updatedService.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(updatedService.UpdatedAt),
	}

	return &response, nil
}

// validatePermissions checks if the updater has permission to update services
func (s *UpdateServiceService) validatePermissions(updaterRole string) error {
	switch updaterRole {
	case common.RoleSuperAdmin, common.RoleAdmin:
		return nil
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}
