package booking

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Cancel struct {
	service bookingService.CancelInterface
}

func NewCancel(service bookingService.CancelInterface) *Cancel {
	return &Cancel{
		service: service,
	}
}

func (h *Cancel) Cancel(c *gin.Context) {
	// Path parameter validation
	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"bookingId": "bookingId 為必填項目",
		})
		return
	}
	parsedBookingID, err := utils.ParseID(bookingID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"bookingId": "bookingId 類型轉換失敗",
		})
		return
	}

	// Input JSON validation
	var req bookingModel.CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim cancelReason
	if req.CancelReason != nil {
		*req.CancelReason = strings.TrimSpace(*req.CancelReason)
	}

	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.Cancel(c.Request.Context(), parsedBookingID, req, customerContext.CustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
