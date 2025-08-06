package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) *Get {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, staffID int64) (*adminStaffModel.GetResponse, error) {
	// Get staff user information
	staffUser, err := s.queries.GetStaffUserByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff user", err)
	}

	// Prepare response with staff information
	response := &adminStaffModel.GetResponse{
		ID:        utils.FormatID(staffUser.ID),
		Username:  staffUser.Username,
		Email:     staffUser.Email,
		Role:      staffUser.Role,
		IsActive:  utils.PgBoolToBool(staffUser.IsActive),
		CreatedAt: utils.PgTimestamptzToTimeString(staffUser.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(staffUser.UpdatedAt),
		Stylist:   nil, // Default to nil
	}

	if staffUser.Role == common.RoleSuperAdmin {
		return response, nil
	}

	// Try to get stylist information
	stylist, err := s.queries.GetStylistByStaffUserID(ctx, staffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist information", err)
	}

	// Convert stylist information
	response.Stylist = &adminStaffModel.GetStaffStylistInfo{
		ID:           utils.FormatID(stylist.ID),
		Name:         utils.PgTextToString(stylist.Name),
		GoodAtShapes: stylist.GoodAtShapes,
		GoodAtColors: stylist.GoodAtColors,
		GoodAtStyles: stylist.GoodAtStyles,
		IsIntrovert:  utils.PgBoolToBool(stylist.IsIntrovert),
		CreatedAt:    utils.PgTimestamptzToTimeString(stylist.CreatedAt),
		UpdatedAt:    utils.PgTimestamptzToTimeString(stylist.UpdatedAt),
	}

	return response, nil
}
