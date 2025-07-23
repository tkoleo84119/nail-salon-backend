package stylist

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateMyStylistService struct {
	queries dbgen.Querier
}

func NewCreateMyStylistService(queries dbgen.Querier) *CreateMyStylistService {
	return &CreateMyStylistService{
		queries: queries,
	}
}

func (s *CreateMyStylistService) CreateMyStylist(ctx context.Context, req stylist.CreateMyStylistRequest, staffUserID int64) (*stylist.CreateMyStylistResponse, error) {
	// Check if stylist already exists for this staff user
	exists, err := s.queries.CheckStylistExistsByStaffUserID(ctx, staffUserID)
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
		StaffUserID:  staffUserID,
		Name:         utils.StringToText(&req.StylistName),
		GoodAtShapes: goodAtShapes,
		GoodAtColors: goodAtColors,
		GoodAtStyles: goodAtStyles,
		IsIntrovert:  utils.BoolPtrToPgBool(&isIntrovert),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create stylist", err)
	}

	// Build response
	response := &stylist.CreateMyStylistResponse{
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
