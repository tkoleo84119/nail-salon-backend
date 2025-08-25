package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineRegister struct {
	service authService.LineRegisterInterface
}

func NewLineRegister(service authService.LineRegisterInterface) *LineRegister {
	return &LineRegister{
		service: service,
	}
}

var validFavoriteShapes = map[string]struct{}{
	"方形": {}, "方圓形": {}, "橢圓形": {}, "圓形": {}, "圓尖形": {}, "尖形": {}, "梯形": {}, "不一定": {},
}

var validFavoriteColors = map[string]struct{}{
	"白色系": {}, "裸色系": {}, "粉色系": {}, "紅色系": {}, "橘色系": {}, "大地色系": {}, "綠色系": {}, "藍色系": {}, "紫色系": {}, "黑色系": {}, "不一定": {},
}

var validFavoriteStyles = map[string]struct{}{
	"暈染": {}, "手繪": {}, "貓眼": {}, "鏡面": {}, "可愛": {}, "法式": {}, "漸層": {}, "氣質溫柔": {}, "個性": {}, "日系": {}, "簡約": {}, "優雅": {}, "典雅": {}, "小眾": {}, "沒有固定": {},
}

func (h *LineRegister) LineRegister(c *gin.Context) {
	var req authModel.LineRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim line idToken, name, email, city, referrer, customerNote
	req.IdToken = strings.TrimSpace(req.IdToken)
	req.Name = strings.TrimSpace(req.Name)
	if req.Email != nil {
		*req.Email = strings.TrimSpace(*req.Email)
	}
	if req.City != nil {
		*req.City = strings.TrimSpace(*req.City)
	}
	if req.Referrer != nil {
		*req.Referrer = strings.TrimSpace(*req.Referrer)
	}
	if req.CustomerNote != nil {
		*req.CustomerNote = strings.TrimSpace(*req.CustomerNote)
	}

	if req.FavoriteShapes != nil {
		for _, shape := range *req.FavoriteShapes {
			if _, ok := validFavoriteShapes[shape]; !ok {
				errorCodes.AbortWithError(c, errorCodes.ValFieldOneof, map[string]string{"favoriteShapes": "favoriteShapes 必須是方形 方圓形 橢圓形 圓形 圓尖形 尖形 梯形 不一定 其中一個值"})
				return
			}
		}
	}

	if req.FavoriteColors != nil {
		for _, color := range *req.FavoriteColors {
			if _, ok := validFavoriteColors[color]; !ok {
				errorCodes.AbortWithError(c, errorCodes.ValFieldOneof, map[string]string{"favoriteColors": "favoriteColors 必須是白色系 裸色系 粉色系 紅色系 橘色系 大地色系 綠色系 藍色系 紫色系 黑色系 不一定 其中一個值"})
				return
			}
		}
	}

	if req.FavoriteStyles != nil {
		for _, style := range *req.FavoriteStyles {
			if _, ok := validFavoriteStyles[style]; !ok {
				errorCodes.AbortWithError(c, errorCodes.ValFieldOneof, map[string]string{"favoriteStyles": "favoriteStyles 必須是暈染 手繪 貓眼 鏡面 可愛 法式 漸層 氣質溫柔 個性 日系 簡約 優雅 典雅 小眾 沒有固定 其中一個值"})
				return
			}
		}
	}

	// Create login context
	loginCtx := authModel.LoginContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Call service
	response, err := h.service.LineRegister(c.Request.Context(), req, loginCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return successful response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
