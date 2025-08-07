package adminStylist

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyStylistHandler struct {
	service adminStylistService.UpdateMyStylistServiceInterface
}

func NewUpdateMyStylistHandler(service adminStylistService.UpdateMyStylistServiceInterface) *UpdateMyStylistHandler {
	return &UpdateMyStylistHandler{
		service: service,
	}
}

var validGoodAtShapes = map[string]struct{}{
	"方形": {}, "方圓形": {}, "橢圓形": {}, "圓形": {}, "圓尖形": {}, "尖形": {}, "梯形": {},
}

var validGoodAtColors = map[string]struct{}{
	"白色系": {}, "裸色系": {}, "粉色系": {}, "紅色系": {}, "橘色系": {}, "大地色系": {}, "綠色系": {}, "藍色系": {}, "紫色系": {}, "黑色系": {},
}

var validGoodAtStyles = map[string]struct{}{
	"暈染": {}, "手繪": {}, "貓眼": {}, "鏡面": {}, "可愛": {}, "法式": {}, "漸層": {}, "氣質溫柔": {}, "個性": {}, "日系": {}, "簡約": {}, "優雅": {}, "典雅": {}, "小眾": {},
}

func (h *UpdateMyStylistHandler) UpdateMyStylist(c *gin.Context) {
	var req adminStylistModel.UpdateMyStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Additional validation: ensure at least one field is provided for update
	if !req.HasUpdate() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	if req.GoodAtShapes != nil {
		for _, shape := range *req.GoodAtShapes {
			if _, ok := validGoodAtShapes[shape]; !ok {
				errorCodes.AbortWithError(c, errorCodes.ValFieldOneof, map[string]string{"goodAtShapes": "goodAtShapes 必須是方形 方圓形 橢圓形 圓形 圓尖形 尖形 梯形 其中一個值"})
				return
			}
		}
	}

	if req.GoodAtColors != nil {
		for _, color := range *req.GoodAtColors {
			if _, ok := validGoodAtColors[color]; !ok {
				errorCodes.AbortWithError(c, errorCodes.ValFieldOneof, map[string]string{"goodAtColors": "goodAtColors 必須是白色系 裸色系 粉色系 紅色系 橘色系 大地色系 綠色系 藍色系 紫色系 黑色系 其中一個值"})
				return
			}
		}
	}

	if req.GoodAtStyles != nil {
		for _, style := range *req.GoodAtStyles {
			if _, ok := validGoodAtStyles[style]; !ok {
				errorCodes.AbortWithError(c, errorCodes.ValFieldOneof, map[string]string{"goodAtStyles": "goodAtStyles 必須是暈染 手繪 貓眼 鏡面 可愛 法式 漸層 氣質溫柔 個性 日系 簡約 優雅 典雅 小眾 其中一個值"})
				return
			}
		}
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.service.UpdateMyStylist(c.Request.Context(), req, staffContext.UserID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
