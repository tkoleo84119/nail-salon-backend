package adminService

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries) *Create {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, req adminServiceModel.CreateRequest, creatorRole string) (*adminServiceModel.CreateResponse, error) {
	// Validate permissions
	if err := s.validatePermissions(creatorRole); err != nil {
		return nil, err
	}

	// Check if service name already exists
	exists, err := s.queries.CheckServiceNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check service existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceAlreadyExists)
	}

	// Generate ID for the new service
	serviceID := utils.GenerateID()

	// Convert price to pgtype.Numeric
	priceNumeric, err := utils.Int64ToPgNumeric(*req.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price", err)
	}

	// Create service
	createdService, err := s.queries.CreateService(ctx, dbgen.CreateServiceParams{
		ID:              serviceID,
		Name:            req.Name,
		Price:           priceNumeric,
		DurationMinutes: *req.DurationMinutes,
		IsAddon:         utils.BoolPtrToPgBool(req.IsAddon),
		IsVisible:       utils.BoolPtrToPgBool(req.IsVisible),
		Note:            utils.StringPtrToPgText(req.Note, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create service", err)
	}

	// Convert to response
	response := &adminServiceModel.CreateResponse{
		ID:              utils.FormatID(createdService.ID),
		Name:            createdService.Name,
		Price:           *req.Price,
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
func (s *Create) validatePermissions(creatorRole string) error {
	switch creatorRole {
	case common.RoleSuperAdmin, common.RoleAdmin:
		return nil
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}
