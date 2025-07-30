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

type UpdateBookingByStaffHandler struct {
	service adminBookingService.UpdateBookingByStaffServiceInterface
}

func NewUpdateBookingByStaffHandler(service adminBookingService.UpdateBookingByStaffServiceInterface) *UpdateBookingByStaffHandler {
	return &UpdateBookingByStaffHandler{
		service: service,
	}
}

func (h *UpdateBookingByStaffHandler) UpdateBookingByStaff(c *gin.Context) {
	// Get path parameters
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"storeId": "storeId 為必填項目"})
		return
	}

	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"bookingId": "bookingId 為必填項目"})
		return
	}

	// Bind request body
	var req adminBookingModel.UpdateBookingByStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.RespondWithEmptyFieldError(c)
		return
	}

	if !req.IsTimeSlotUpdateComplete() {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"request": "timeSlotId、mainServiceId、subServiceIds 必須一起傳入",
		})
		return
	}

	// Call service
	response, err := h.service.UpdateBookingByStaff(c.Request.Context(), storeID, bookingID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
