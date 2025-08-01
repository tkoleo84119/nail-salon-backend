package adminAuth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StaffRefreshTokenService struct {
	repo      *sqlxRepo.Repositories
	jwtConfig config.JWTConfig
}

func NewStaffRefreshTokenService(repo *sqlxRepo.Repositories, jwtConfig config.JWTConfig) *StaffRefreshTokenService {
	return &StaffRefreshTokenService{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

func (s *StaffRefreshTokenService) StaffRefreshToken(ctx context.Context, req adminAuthModel.StaffRefreshTokenRequest) (*adminAuthModel.StaffRefreshTokenResponse, error) {
	// Validate refresh token exists and is not revoked/expired
	refreshTokenInfo, err := s.repo.StaffUserTokens.GetStaffUserTokenValid(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthRefreshTokenInvalid)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to validate refresh token", err)
	}

	// Get staff user information to rebuild JWT claims
	staffUser, err := s.repo.Staff.GetStaffUserByID(ctx, refreshTokenInfo.StaffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get staff user", err)
	}

	// Get store access for the staff user
	storeList, err := s.getStoreAccess(ctx, staffUser)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store access", err)
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(s.jwtConfig, staffUser.ID, staffUser.Username, staffUser.Role, storeList)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	// Build response
	return &adminAuthModel.StaffRefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   s.jwtConfig.ExpiryHours * 3600,
		User: adminAuthModel.User{
			ID:        utils.FormatID(staffUser.ID),
			Username:  staffUser.Username,
			Role:      staffUser.Role,
			StoreList: storeList,
		},
	}, nil
}

// getStoreAccess is a helper method to get store access based on staff role
func (s *StaffRefreshTokenService) getStoreAccess(ctx context.Context, staffUser *sqlxRepo.GetStaffUserByIDResponse) ([]common.Store, error) {
	var storeList []common.Store

	switch staffUser.Role {
	case common.RoleSuperAdmin:
		// Super admin has access to all active stores
		stores, err := s.repo.Store.GetAllStore(ctx, nil)
		if err != nil {
			return nil, err
		}
		for _, store := range stores {
			storeList = append(storeList, common.Store{
				ID:   utils.FormatID(store.ID),
				Name: store.Name,
			})
		}
	case common.RoleAdmin, common.RoleManager, common.RoleStylist:
		// Other roles have access based on store access table
		storeAccess, err := s.repo.StaffUserStoreAccess.GetStaffUserStoreAccessByStaffId(ctx, staffUser.ID, nil)
		if err != nil {
			return nil, err
		}
		for _, access := range storeAccess {
			storeList = append(storeList, common.Store{
				ID:   utils.FormatID(access.StoreID),
				Name: access.Name,
			})
		}
	}

	return storeList, nil
}
