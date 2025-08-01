package adminAuth

import (
	"context"
	"net/netip"
	"time"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StaffLoginService struct {
	repo      *sqlxRepo.Repositories
	jwtConfig config.JWTConfig
}

func NewStaffLoginService(repo *sqlxRepo.Repositories, jwtConfig config.JWTConfig) *StaffLoginService {
	return &StaffLoginService{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

func (s *StaffLoginService) StaffLogin(ctx context.Context, req adminAuthModel.StaffLoginRequest, loginCtx adminAuthModel.StaffLoginContext) (*adminAuthModel.StaffLoginResponse, error) {
	staffUser, err := s.repo.Staff.GetStaffUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthInvalidCredentials)
	}
	if !utils.CheckPassword(req.Password, staffUser.PasswordHash) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthInvalidCredentials)
	}
	if !staffUser.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthInvalidCredentials)
	}

	// Get store access based on role
	storeList, err := s.getStoreAccess(ctx, staffUser)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store access", err)
	}

	// Generate tokens
	accessToken, err := utils.GenerateJWT(s.jwtConfig, staffUser.ID, staffUser.Username, staffUser.Role, storeList)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// Store refresh token
	tokenInfo := adminAuthModel.StaffTokenInfo{
		StaffUserID:  staffUser.ID,
		RefreshToken: refreshToken,
		Context:      loginCtx,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.storeRefreshToken(ctx, tokenInfo); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to store refresh token", err)
	}

	response := &adminAuthModel.StaffLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtConfig.ExpiryHours * 3600,
		User: adminAuthModel.User{
			ID:        utils.FormatID(staffUser.ID),
			Username:  staffUser.Username,
			Role:      staffUser.Role,
			StoreList: storeList,
		},
	}

	return response, nil
}

// getStoreAccess retrieves store access based on user role
func (s *StaffLoginService) getStoreAccess(ctx context.Context, staffUser *sqlxRepo.GetStaffUserByUsernameResponse) ([]common.Store, error) {
	var storeList []common.Store

	// SUPER_ADMIN can access all stores
	if staffUser.Role == common.RoleSuperAdmin {
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
	} else {
		// Get specific store access for other roles
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

// storeRefreshToken stores the refresh token in database
func (s *StaffLoginService) storeRefreshToken(ctx context.Context, tokenInfo adminAuthModel.StaffTokenInfo) error {
	// Convert IP address to string pointer for PostgreSQL INET type
	var ipAddr *string
	if tokenInfo.Context.IPAddress != "" {
		if _, err := netip.ParseAddr(tokenInfo.Context.IPAddress); err == nil {
			ipAddr = &tokenInfo.Context.IPAddress
		}
	}

	_, err := s.repo.StaffUserTokens.CreateStaffUserToken(ctx, sqlxRepo.CreateStaffUserTokenParams{
		ID:           utils.GenerateID(),
		StaffUserID:  tokenInfo.StaffUserID,
		RefreshToken: tokenInfo.RefreshToken,
		UserAgent:    utils.StringPtrToPgText(&tokenInfo.Context.UserAgent, true),
		IpAddress:    ipAddr,
		ExpiredAt:    utils.TimeToPgTimestamptz(tokenInfo.ExpiresAt),
	})

	return err
}
