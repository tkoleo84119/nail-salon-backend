package adminBooking

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminBookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CancelBookingHandler struct {
	service adminBookingService.CancelBookingServiceInterface
}

func NewCancelBookingHandler(service adminBookingService.CancelBookingServiceInterface) *CancelBookingHandler {
	return &CancelBookingHandler{
		service: service,
	}
}

func (h *CancelBookingHandler) CancelBooking(c *gin.Context) {
	// Get path parameters
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{"storeId": "storeId 為必填項目"})
		return
	}

	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{"bookingId": "bookingId 為必填項目"})
		return
	}

	// Parse request body
	var req adminBookingModel.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Call service
	result, err := h.service.CancelBooking(c.Request.Context(), storeID, bookingID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(result))
}
