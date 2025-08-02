package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMyStaffService struct {
	repo *sqlxRepo.Repositories
}

func NewGetMyStaffService(repo *sqlxRepo.Repositories) *GetMyStaffService {
	return &GetMyStaffService{
		repo: repo,
	}
}

func (s *GetMyStaffService) GetMyStaff(ctx context.Context, staffUserID int64) (*adminStaffModel.GetMyStaffResponse, error) {
	// Get staff user information
	staffUser, err := s.repo.Staff.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff user", err)
	}

	// Convert to response format
	response := &adminStaffModel.GetMyStaffResponse{
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

	stylist, err := s.repo.Stylist.GetStylistByStaffUserID(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist", err)
	}

	response.Stylist = &adminStaffModel.StaffStylistInfo{
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
