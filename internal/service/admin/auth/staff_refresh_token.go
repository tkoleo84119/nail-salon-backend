package adminAuth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StaffRefreshTokenService struct {
	queries   *dbgen.Queries
	jwtConfig config.JWTConfig
}

func NewStaffRefreshTokenService(queries *dbgen.Queries, jwtConfig config.JWTConfig) *StaffRefreshTokenService {
	return &StaffRefreshTokenService{
		queries:   queries,
		jwtConfig: jwtConfig,
	}
}

func (s *StaffRefreshTokenService) StaffRefreshToken(ctx context.Context, req adminAuthModel.StaffRefreshTokenRequest) (*adminAuthModel.StaffRefreshTokenResponse, error) {
	// Validate refresh token exists and is not revoked/expired
	tokenRecord, err := s.queries.GetValidStaffUserToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthRefreshTokenInvalid)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to validate refresh token", err)
	}

	// Get staff user information to rebuild JWT claims
	staffUser, err := s.queries.GetStaffUserByID(ctx, tokenRecord.StaffUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthRefreshTokenInvalid)
		}
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
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to generate access token", err)
	}

	// Build response
	return &adminAuthModel.StaffRefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   s.jwtConfig.ExpiryHours * 3600,
	}, nil
}

// getStoreAccess is a helper method to get store access based on staff role
func (s *StaffRefreshTokenService) getStoreAccess(ctx context.Context, staffUser dbgen.StaffUser) ([]common.Store, error) {
	var storeList []common.Store

	switch staffUser.Role {
	case "SUPER_ADMIN":
		// Super admin has access to all active stores
		stores, err := s.queries.GetAllActiveStores(ctx)
		if err != nil {
			return nil, err
		}
		for _, store := range stores {
			storeList = append(storeList, common.Store{
				ID:   utils.FormatID(store.ID),
				Name: store.Name,
			})
		}
	case "ADMIN", "MANAGER", "STYLIST":
		// Other roles have access based on store access table
		storeAccess, err := s.queries.GetStaffUserStoreAccess(ctx, staffUser.ID)
		if err != nil {
			return nil, err
		}
		for _, access := range storeAccess {
			storeList = append(storeList, common.Store{
				ID:   utils.FormatID(access.StoreID),
				Name: access.StoreName,
			})
		}
	}

	return storeList, nil
}