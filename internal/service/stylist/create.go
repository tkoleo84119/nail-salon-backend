package stylist

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStylistService struct {
	queries dbgen.Querier
}

func NewCreateStylistService(queries dbgen.Querier) *CreateStylistService {
	return &CreateStylistService{
		queries: queries,
	}
}

func (s *CreateStylistService) CreateStylist(ctx context.Context, req stylist.CreateStylistRequest, staffUserID int64) (*stylist.CreateStylistResponse, error) {
	// Get staff user info to check role
	staffUser, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "failed to get staff user", err)
	}

	// Check if user is SUPER_ADMIN (not allowed to create stylist)
	if staffUser.Role == staff.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if stylist already exists for this staff user
	exists, err := s.queries.CheckStylistExistsByStaffUserID(ctx, pgtype.Int8{Int64: staffUserID, Valid: true})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check stylist existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistAlreadyExists)
	}

	// Generate ID for new stylist
	stylistID := utils.GenerateID()

	// Handle default value for IsIntrovert
	isIntrovert := false
	if req.IsIntrovert != nil {
		isIntrovert = *req.IsIntrovert
	}

	// Handle empty slices
	goodAtShapes := req.GoodAtShapes
	if goodAtShapes == nil {
		goodAtShapes = []string{}
	}
	goodAtColors := req.GoodAtColors
	if goodAtColors == nil {
		goodAtColors = []string{}
	}
	goodAtStyles := req.GoodAtStyles
	if goodAtStyles == nil {
		goodAtStyles = []string{}
	}

	// Create stylist record
	createdStylist, err := s.queries.CreateStylist(ctx, dbgen.CreateStylistParams{
		ID:           stylistID,
		StaffUserID:  pgtype.Int8{Int64: staffUserID, Valid: true},
		Name:         pgtype.Text{String: req.StylistName, Valid: true},
		GoodAtShapes: goodAtShapes,
		GoodAtColors: goodAtColors,
		GoodAtStyles: goodAtStyles,
		IsIntrovert:  pgtype.Bool{Bool: isIntrovert, Valid: true},
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create stylist", err)
	}

	// Build response
	response := &stylist.CreateStylistResponse{
		ID:           utils.FormatID(createdStylist.ID),
		StaffUserID:  utils.FormatID(staffUserID),
		StylistName:  createdStylist.Name.String,
		GoodAtShapes: createdStylist.GoodAtShapes,
		GoodAtColors: createdStylist.GoodAtColors,
		GoodAtStyles: createdStylist.GoodAtStyles,
		IsIntrovert:  createdStylist.IsIntrovert.Bool,
	}

	return response, nil
}
