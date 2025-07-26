package adminStaff

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type GetStaffListService struct {
	repository sqlx.StaffUserRepositoryInterface
}

func NewGetStaffListService(repository sqlx.StaffUserRepositoryInterface) *GetStaffListService {
	return &GetStaffListService{
		repository: repository,
	}
}

func (s *GetStaffListService) GetStaffList(ctx context.Context, req adminStaffModel.GetStaffListRequest) (*adminStaffModel.GetStaffListResponse, error) {
	// Delegate to repository layer which handles all filtering, pagination, and response formatting
	response, err := s.repository.GetStaffList(ctx, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get staff list", err)
	}

	return response, nil
}