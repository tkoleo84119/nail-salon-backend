package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyBookingHandler struct {
	service bookingService.UpdateMyBookingServiceInterface
}

func NewUpdateMyBookingHandler(service bookingService.UpdateMyBookingServiceInterface) *UpdateMyBookingHandler {
	return &UpdateMyBookingHandler{
		service: service,
	}
}

func (h *UpdateMyBookingHandler) UpdateMyBooking(c *gin.Context) {
	// Input JSON validation
	var req bookingModel.UpdateMyBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Path parameter validation
	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"bookingId": "bookingId為必填項目",
		})
		return
	}

	// Business logic validation - HasUpdates check
	if !req.HasUpdates() {
		errorCodes.RespondWithEmptyFieldError(c)
		return
	}

	// Business logic validation - Time slot update completeness
	if !req.IsTimeSlotUpdateComplete() {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"request": "storeId、stylistId、timeSlotId、mainServiceId、subServiceIds 必須一起傳入",
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
	response, err := h.service.UpdateMyBooking(c.Request.Context(), bookingID, req, *customerContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
