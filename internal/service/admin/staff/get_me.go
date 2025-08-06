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

type GetMe struct {
	queries *dbgen.Queries
}

func NewGetMe(queries *dbgen.Queries) *GetMe {
	return &GetMe{
		queries: queries,
	}
}

func (s *GetMe) GetMe(ctx context.Context, staffUserID int64) (*adminStaffModel.GetMeResponse, error) {
	// Get staff user information
	staffUser, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff user", err)
	}

	// Convert to response format
	response := &adminStaffModel.GetMeResponse{
		ID:        utils.FormatID(staffUser.ID),
		Username:  staffUser.Username,
		Email:     staffUser.Email,
		Role:      staffUser.Role,
		IsActive:  utils.PgBoolToBool(staffUser.IsActive),
		CreatedAt: utils.PgTimestamptzToTimeString(staffUser.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(staffUser.UpdatedAt),
		Stylist:   nil, // default is nil
	}

	if staffUser.Role == common.RoleSuperAdmin {
		return response, nil
	}

	stylist, err := s.queries.GetStylistByStaffUserID(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist", err)
	}

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
