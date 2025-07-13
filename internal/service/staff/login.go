package staff

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LoginService struct {
	queries   dbgen.Querier
	jwtConfig config.JWTConfig
}

func NewLoginService(queries dbgen.Querier, jwtConfig config.JWTConfig) *LoginService {
	return &LoginService{
		queries:   queries,
		jwtConfig: jwtConfig,
	}
}

func (s *LoginService) Login(ctx context.Context, req staff.LoginRequest, loginCtx staff.LoginContext) (*staff.LoginResponse, error) {
	staffUser, err := s.queries.GetStaffUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, staffUser.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Get store access based on role
	storeList, err := s.getStoreAccess(ctx, staffUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get store access: %w", err)
	}

	// Generate tokens
	accessToken, err := utils.GenerateJWT(s.jwtConfig, staffUser.ID, staffUser.Username, staffUser.Role, storeList)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token
	tokenInfo := staff.TokenInfo{
		StaffUserID:  staffUser.ID,
		RefreshToken: refreshToken,
		Context:      loginCtx,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.storeRefreshToken(ctx, tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	response := &staff.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtConfig.ExpiryHours * 3600,
		User: staff.User{
			ID:        strconv.FormatInt(staffUser.ID, 10),
			Username:  staffUser.Username,
			Role:      staffUser.Role,
			StoreList: storeList,
		},
	}

	return response, nil
}

// getStoreAccess retrieves store access based on user role
func (s *LoginService) getStoreAccess(ctx context.Context, staffUser dbgen.StaffUser) ([]utils.Store, error) {
	var storeList []utils.Store

	// SUPER_ADMIN can access all stores
	if staffUser.Role == staff.RoleSuperAdmin {
		stores, err := s.queries.GetAllActiveStores(ctx)
		if err != nil {
			return nil, err
		}
		for _, store := range stores {
			storeList = append(storeList, utils.Store{
				ID:   store.ID,
				Name: store.Name,
			})
		}
	} else {
		// Get specific store access for other roles
		storeAccess, err := s.queries.GetStaffUserStoreAccess(ctx, staffUser.ID)
		if err != nil {
			return nil, err
		}
		for _, access := range storeAccess {
			storeList = append(storeList, utils.Store{
				ID:   access.StoreID,
				Name: access.StoreName,
			})
		}
	}

	return storeList, nil
}

// storeRefreshToken stores the refresh token in database
func (s *LoginService) storeRefreshToken(ctx context.Context, tokenInfo staff.TokenInfo) error {
	// Parse IP address
	var ipAddr *netip.Addr
	if addr, err := netip.ParseAddr(tokenInfo.Context.IPAddress); err == nil {
		ipAddr = &addr
	}

	tokenID := utils.GenerateID()

	_, err := s.queries.CreateStaffUserToken(ctx, dbgen.CreateStaffUserTokenParams{
		ID:           tokenID,
		StaffUserID:  tokenInfo.StaffUserID,
		RefreshToken: tokenInfo.RefreshToken,
		UserAgent:    pgtype.Text{String: tokenInfo.Context.UserAgent, Valid: tokenInfo.Context.UserAgent != ""},
		IpAddress:    ipAddr,
		ExpiredAt:    pgtype.Timestamptz{Time: tokenInfo.ExpiresAt, Valid: true},
	})

	return err
}
