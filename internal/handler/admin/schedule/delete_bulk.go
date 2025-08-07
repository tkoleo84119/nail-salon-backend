package adminSchedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteBulk struct {
	service adminScheduleService.DeleteBulkInterface
}

func NewDeleteBulk(service adminScheduleService.DeleteBulkInterface) *DeleteBulk {
	return &DeleteBulk{
		service: service,
	}
}

func (h *DeleteBulk) DeleteBulk(c *gin.Context) {
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId 為必填項目",
		})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	// Parse and validate request
	var req adminScheduleModel.DeleteBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	parsedStylistID, err := utils.ParseID(req.StylistID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"stylistId": "stylistId 類型轉換失敗",
		})
		return
	}

	parsedScheduleIDs := make([]int64, len(req.ScheduleIDs))
	for i, scheduleID := range req.ScheduleIDs {
		parsedScheduleID, err := utils.ParseID(scheduleID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"scheduleId": "scheduleId 類型轉換失敗",
			})
		}
		parsedScheduleIDs[i] = parsedScheduleID
	}

	parsedReq := adminScheduleModel.DeleteBulkParsedRequest{
		StylistID:   parsedStylistID,
		ScheduleIDs: parsedScheduleIDs,
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	updaterStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		updaterStoreIDs[i] = store.ID
	}

	// Call service
	response, err := h.service.DeleteBulk(c.Request.Context(), parsedStoreID, parsedReq, staffContext.UserID, staffContext.Role, updaterStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response with 204 No Content but with data
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
