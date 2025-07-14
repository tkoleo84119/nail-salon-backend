package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
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
			errors := map[string]string{"token": "access_token 缺失"}
			c.JSON(http.StatusUnauthorized, common.ErrorResponse("認證失敗", errors))
			c.Abort()
			return
		}

		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errors := map[string]string{"token": "access_token 格式錯誤"}
			c.JSON(http.StatusUnauthorized, common.ErrorResponse("認證失敗", errors))
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := utils.ValidateJWT(cfg.JWT, token)
		if err != nil {
			errors := map[string]string{"token": "access_token 無效或已過期"}
			c.JSON(http.StatusUnauthorized, common.ErrorResponse("認證失敗", errors))
			c.Abort()
			return
		}

		if err := validateStaffToken(c, db, claims); err != nil {
			errors := map[string]string{"token": "員工認證失敗"}
			c.JSON(http.StatusUnauthorized, common.ErrorResponse("認證失敗", errors))
			c.Abort()
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
			errors := map[string]string{"auth": "未找到使用者認證資訊"}
			c.JSON(http.StatusUnauthorized, common.ErrorResponse("認證失敗", errors))
			c.Abort()
			return
		}

		if !hasRequiredRole(staffContext.Role, allowedRoles) {
			errors := map[string]string{"permission": "權限不足"}
			c.JSON(http.StatusForbidden, common.ErrorResponse("權限不足", errors))
			c.Abort()
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
