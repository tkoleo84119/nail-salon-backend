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

type Create struct {
	service adminBookingService.CreateInterface
}

func NewCreate(service adminBookingService.CreateInterface) *Create {
	return &Create{service: service}
}

func (h *Create) Create(c *gin.Context) {
	// Get path parameter
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"storeId": "storeId 為必填項目"})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"storeId": "storeId 類型轉換失敗"})
		return
	}

	// Parse request body
	var req adminBookingModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim note
	if req.StoreNote != nil {
		*req.StoreNote = strings.TrimSpace(*req.StoreNote)
	}

	parsedCustomerID, err := utils.ParseID(req.CustomerID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"customerId": "customerId 類型轉換失敗"})
		return
	}

	parsedTimeSlotID, err := utils.ParseID(req.TimeSlotID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"timeSlotId": "timeSlotId 類型轉換失敗"})
		return
	}

	parsedMainServiceID, err := utils.ParseID(req.MainServiceID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"mainServiceId": "mainServiceId 類型轉換失敗"})
		return
	}

	parsedSubServiceIDs := []int64{}
	for _, subServiceID := range *req.SubServiceIDs {
		parsedSubServiceID, err := utils.ParseID(subServiceID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"subServiceIds": "subServiceIds 類型轉換失敗"})
			return
		}
		parsedSubServiceIDs = append(parsedSubServiceIDs, parsedSubServiceID)
	}

	isChatEnabled := false
	if req.IsChatEnabled != nil {
		isChatEnabled = *req.IsChatEnabled
	}

	parsedRequest := adminBookingModel.CreateParsedRequest{
		CustomerID:    parsedCustomerID,
		TimeSlotID:    parsedTimeSlotID,
		MainServiceID: parsedMainServiceID,
		SubServiceIDs: parsedSubServiceIDs,
		IsChatEnabled: isChatEnabled,
		StoreNote:     req.StoreNote,
	}

	// Get staff context from JWT middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	storeIds := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeIds[i] = store.ID
	}

	// Call service
	booking, err := h.service.Create(c.Request.Context(), parsedStoreID, parsedRequest, staffContext.Role, storeIds, staffContext.UserID, staffContext.Username)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response with 201 Created
	c.JSON(http.StatusCreated, common.SuccessResponse(booking))
}
