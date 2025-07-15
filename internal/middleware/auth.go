package middleware

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

const (
	UserContextKey = "user"
)

func JWTAuth(cfg config.Config, db dbgen.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errorCodes.AbortWithError(c, errorCodes.AuthTokenMissing, nil)
			return
		}

		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errorCodes.AbortWithError(c, errorCodes.AuthTokenFormatError, nil)
			return
		}

		token := tokenParts[1]
		claims, err := utils.ValidateJWT(cfg.JWT, token)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.AuthTokenInvalid, nil)
			return
		}

		if err := validateStaffToken(c, db, claims); err != nil {
			errorCodes.AbortWithError(c, errorCodes.AuthStaffFailed, nil)
			return
		}

		c.Next()
	}
}

// validate staff is active with token claims
func validateStaffToken(c *gin.Context, db dbgen.Querier, claims *common.JWTClaims) error {
	staff, err := db.GetStaffUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		return err
	}

	if !staff.IsActive.Bool {
		return err
	}

	c.Set(UserContextKey, claims.StaffContext)
	return nil
}

func GetStaffFromContext(c *gin.Context) (*common.StaffContext, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}

	staffContext, ok := user.(common.StaffContext)
	return &staffContext, ok
}

// check if staff has required role
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		staffContext, exists := GetStaffFromContext(c)
		if !exists {
			errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
			return
		}

		if !hasRequiredRole(staffContext.Role, allowedRoles) {
			errorCodes.AbortWithError(c, errorCodes.AuthPermissionDenied, nil)
			return
		}

		c.Next()
	}
}

func hasRequiredRole(userRole string, allowedRoles []string) bool {
	return slices.Contains(allowedRoles, userRole)
}

func RequireAdminRoles() gin.HandlerFunc {
	return RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin)
}

func RequireManagerOrAbove() gin.HandlerFunc {
	return RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin, staff.RoleManager)
}

func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRoles(staff.RoleSuperAdmin)
}

func RequireAnyStaffRole() gin.HandlerFunc {
	return RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin, staff.RoleManager, staff.RoleStylist)
}
