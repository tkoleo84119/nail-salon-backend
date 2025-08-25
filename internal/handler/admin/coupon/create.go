package adminCoupon

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminCouponService.CreateInterface
}

func NewCreate(service adminCouponService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	var req adminCouponModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name, displayName, code, note
	req.Name = strings.TrimSpace(req.Name)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.Code = strings.TrimSpace(req.Code)
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	response, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
