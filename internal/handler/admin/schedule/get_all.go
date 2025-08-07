package adminSchedule

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminScheduleService.GetAllInterface
}

func NewGetAll(service adminScheduleService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Get store ID from path parameter
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

	// Parse and validate query parameters
	var req adminScheduleModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	startDate, err := utils.DateStringToTime(req.StartDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{"startDate": "startDate 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD)"})
		return
	}
	endDate, err := utils.DateStringToTime(req.EndDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{"endDate": "endDate 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD)"})
		return
	}

	var parsedStylistIDList []int64
	if req.StylistID != nil && *req.StylistID != "" {
		stringStylistIDList := strings.Split(*req.StylistID, ",")
		for _, stylistID := range stringStylistIDList {
			parsedStylistID, err := utils.ParseID(stylistID)
			if err != nil {
				errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"stylistId": "stylistId 類型轉換失敗"})
				return
			}
			parsedStylistIDList = append(parsedStylistIDList, parsedStylistID)
		}
	}

	parsedReq := adminScheduleModel.GetAllParsedRequest{
		StylistID:   &parsedStylistIDList,
		StartDate:   startDate,
		EndDate:     endDate,
		IsAvailable: req.IsAvailable,
	}

	// Get staff context from JWT middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	storeIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeIDs[i] = store.ID
	}

	// Service layer call
	response, err := h.service.GetAll(c.Request.Context(), parsedStoreID, parsedReq, staffContext.Role, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
