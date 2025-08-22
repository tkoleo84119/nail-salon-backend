package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMe struct {
	service customerService.UpdateMeInterface
}

func NewUpdateMe(service customerService.UpdateMeInterface) *UpdateMe {
	return &UpdateMe{
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

// UpdateMyCustomer handles PATCH /api/customers/me
func (h *UpdateMe) UpdateMe(c *gin.Context) {
	var req customerModel.UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
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

	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	customerID := customerContext.CustomerID
	result, err := h.service.UpdateMe(c.Request.Context(), customerID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(result))
}
