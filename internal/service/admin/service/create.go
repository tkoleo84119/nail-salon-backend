package adminService

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateServiceService struct {
	repo *sqlxRepo.Repositories
}

func NewCreateServiceService(repo *sqlxRepo.Repositories) *CreateServiceService {
	return &CreateServiceService{
		repo: repo,
	}
}

func (s *CreateServiceService) CreateService(ctx context.Context, req adminServiceModel.CreateServiceRequest, creatorRole string) (*adminServiceModel.CreateServiceResponse, error) {
	// Validate permissions
	if err := s.validatePermissions(creatorRole); err != nil {
		return nil, err
	}

	// Check if service name already exists
	exists, err := s.repo.Service.CheckServiceNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check service existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceAlreadyExists)
	}

	// Generate ID for the new service
	serviceID := utils.GenerateID()

	// Convert price to pgtype.Numeric
	priceNumeric, err := utils.Int64ToPgNumeric(req.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to convert price", err)
	}

	// Create service
	createdService, err := s.repo.Service.CreateService(ctx, sqlxRepo.CreateServiceParams{
		ID:              serviceID,
		Name:            req.Name,
		Price:           priceNumeric,
		DurationMinutes: req.DurationMinutes,
		IsAddon:         utils.BoolPtrToPgBool(&req.IsAddon),
		IsVisible:       utils.BoolPtrToPgBool(&req.IsVisible),
		Note:            utils.StringPtrToPgText(&req.Note, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create service", err)
	}

	// Convert to response
	response := &adminServiceModel.CreateServiceResponse{
		ID:              utils.FormatID(createdService.ID),
		Name:            createdService.Name,
		Price:           req.Price,
		DurationMinutes: createdService.DurationMinutes,
		IsAddon:         utils.PgBoolToBool(createdService.IsAddon),
		IsVisible:       utils.PgBoolToBool(createdService.IsVisible),
		IsActive:        utils.PgBoolToBool(createdService.IsActive),
		Note:            utils.PgTextToString(createdService.Note),
		CreatedAt:       utils.PgTimestamptzToTimeString(createdService.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(createdService.UpdatedAt),
	}

	return response, nil
}

// validatePermissions checks if the creator has permission to create services
func (s *CreateServiceService) validatePermissions(creatorRole string) error {
	switch creatorRole {
	case common.RoleSuperAdmin, common.RoleAdmin:
		return nil
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}
