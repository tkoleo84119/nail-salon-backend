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

type Update struct {
	service bookingService.UpdateInterface
}

func NewUpdate(service bookingService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
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

	// Input JSON validation
	var req bookingModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	storeId, err := utils.ParseID(*req.StoreId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	stylistId, err := utils.ParseID(*req.StylistId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"stylistId": "stylistId 類型轉換失敗",
		})
		return
	}

	timeSlotId, err := utils.ParseID(*req.TimeSlotId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"timeSlotId": "timeSlotId 類型轉換失敗",
		})
		return
	}

	mainServiceId, err := utils.ParseID(*req.MainServiceId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"mainServiceId": "mainServiceId 類型轉換失敗",
		})
		return
	}

	subServiceIds := make([]int64, len(*req.SubServiceIds))
	for i, subServiceId := range *req.SubServiceIds {
		subServiceId, err := utils.ParseID(subServiceId)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"subServiceId": "subServiceId 類型轉換失敗",
			})
			return
		}
		subServiceIds[i] = subServiceId
	}

	parsedRequest := bookingModel.UpdateParsedRequest{
		StoreId:       &storeId,
		StylistId:     &stylistId,
		TimeSlotId:    &timeSlotId,
		MainServiceId: &mainServiceId,
		SubServiceIds: &subServiceIds,
		Note:          req.Note,
		IsChatEnabled: req.IsChatEnabled,
	}

	// Business logic validation - HasUpdates check
	if !parsedRequest.HasUpdates() {
		errorCodes.RespondWithEmptyFieldError(c)
		return
	}

	// Business logic validation - Time slot update completeness
	if !parsedRequest.IsTimeSlotUpdateComplete() {
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
	response, err := h.service.Update(c.Request.Context(), parsedBookingID, parsedRequest, customerContext.CustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
