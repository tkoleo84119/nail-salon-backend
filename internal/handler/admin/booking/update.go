package adminBooking

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminBookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminBookingService.UpdateInterface
}

func NewUpdate(service adminBookingService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	// Get path parameters
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"storeId": "storeId 為必填項目"})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"bookingId": "bookingId 為必填項目"})
		return
	}
	parsedBookingID, err := utils.ParseID(bookingID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"bookingId": "bookingId 類型轉換失敗",
		})
		return
	}

	// Bind request body
	var req adminBookingModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	if !req.IsTimeSlotUpdateComplete() {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"request": "stylistId、timeSlotId、mainServiceId、subServiceIds 必須一起傳入",
		})
		return
	}

	// trim note
	if req.StoreNote != nil {
		*req.StoreNote = strings.TrimSpace(*req.StoreNote)
	}

	var stylistId *int64
	if req.StylistID != nil {
		parsed, err := utils.ParseID(*req.StylistID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"stylistId": "stylistId 類型轉換失敗",
			})
		}
		stylistId = &parsed
	}

	var timeSlotId *int64
	if req.TimeSlotID != nil {
		parsed, err := utils.ParseID(*req.TimeSlotID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"timeSlotId": "timeSlotId 類型轉換失敗",
			})
		}
		timeSlotId = &parsed
	}

	var mainServiceId *int64
	if req.MainServiceID != nil {
		parsed, err := utils.ParseID(*req.MainServiceID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"mainServiceId": "mainServiceId 類型轉換失敗",
			})
		}
		mainServiceId = &parsed
	}

	var subServiceIds []int64
	if req.SubServiceIDs != nil {
		for _, subServiceId := range *req.SubServiceIDs {
			subServiceIdInt, err := utils.ParseID(subServiceId)
			if err != nil {
				errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
					"subServiceIds": "subServiceIds 類型轉換失敗",
				})
			}
			subServiceIds = append(subServiceIds, subServiceIdInt)
		}
	}

	parsedReq := adminBookingModel.UpdateParsedRequest{
		StylistID:     stylistId,
		TimeSlotID:    timeSlotId,
		MainServiceID: mainServiceId,
		SubServiceIDs: subServiceIds,
		IsChatEnabled: req.IsChatEnabled,
		StoreNote:     req.StoreNote,
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	response, err := h.service.Update(c.Request.Context(), parsedStoreID, parsedBookingID, parsedReq, staffContext.Username)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
