package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	service bookingService.GetInterface
}

func NewGet(service bookingService.GetInterface) *Get {
	return &Get{
		service: service,
	}
}

func (h *Get) Get(c *gin.Context) {
	// Path parameter validation
	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"bookingId": "bookingId為必填項目",
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

	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.Get(c.Request.Context(), parsedBookingID, customerContext.CustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
