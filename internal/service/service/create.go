package service

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateServiceService struct {
	queries dbgen.Querier
}

func NewCreateServiceService(queries dbgen.Querier) *CreateServiceService {
	return &CreateServiceService{
		queries: queries,
	}
}

func (s *CreateServiceService) CreateService(ctx context.Context, req service.CreateServiceRequest, creatorRole string) (*service.CreateServiceResponse, error) {
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
	priceNumeric, err := utils.Int64ToPgNumeric(req.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to convert price", err)
	}

	// Convert note to pgtype.Text
	noteText := utils.StringToText(&req.Note)

	// Create service
	createdService, err := s.queries.CreateService(ctx, dbgen.CreateServiceParams{
		ID:              serviceID,
		Name:            req.Name,
		Price:           priceNumeric,
		DurationMinutes: req.DurationMinutes,
		IsAddon:         utils.BoolPtrToPgBool(&req.IsAddon),
		IsVisible:       utils.BoolPtrToPgBool(&req.IsVisible),
		Note:            noteText,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create service", err)
	}

	// Convert to response
	response := &service.CreateServiceResponse{
		ID:              utils.FormatID(createdService.ID),
		Name:            createdService.Name,
		Price:           req.Price,
		DurationMinutes: createdService.DurationMinutes,
		IsAddon:         createdService.IsAddon.Bool,
		IsVisible:       createdService.IsVisible.Bool,
		IsActive:        createdService.IsActive.Bool,
		Note:            createdService.Note.String,
	}

	return response, nil
}

// validatePermissions checks if the creator has permission to create services
func (s *CreateServiceService) validatePermissions(creatorRole string) error {
	switch creatorRole {
	case staff.RoleSuperAdmin, staff.RoleAdmin:
		return nil
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}
