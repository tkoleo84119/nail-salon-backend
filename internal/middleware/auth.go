package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
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
