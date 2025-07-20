package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

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
	checkService, err := s.queries.GetServiceByName(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check service existence", err)
	}
	if checkService.Name != "" {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNameAlreadyExists)
	}

	// Generate ID for the new service
	serviceID := utils.GenerateID()

	// Convert price to pgtype.Numeric
	priceNumeric := pgtype.Numeric{}
	err = priceNumeric.Scan(req.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to convert price", err)
	}

	// Convert note to pgtype.Text
	noteText := pgtype.Text{}
	if req.Note != "" {
		err = noteText.Scan(req.Note)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to convert note", err)
		}
	}

	// Create service
	createdService, err := s.queries.CreateService(ctx, dbgen.CreateServiceParams{
		ID:              serviceID,
		Name:            req.Name,
		Price:           priceNumeric,
		DurationMinutes: req.DurationMinutes,
		IsAddon:         pgtype.Bool{Bool: req.IsAddon, Valid: true},
		IsVisible:       pgtype.Bool{Bool: req.IsVisible, Valid: true},
		IsActive:        pgtype.Bool{Bool: true, Valid: true}, // Default to true
		Note:            noteText,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create service", err)
	}

	// Convert to response
	response := &service.CreateServiceResponse{
		ID:              utils.FormatID(createdService.ID),
		Name:            createdService.Name,
		Price:           req.Price, // Use original request price as int64
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
