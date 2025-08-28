package adminReport

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminReportModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/report"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminReportService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/report"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetPerformanceMe struct {
	service adminReportService.GetPerformanceMeInterface
}

func NewGetPerformanceMe(service adminReportService.GetPerformanceMeInterface) *GetPerformanceMe {
	return &GetPerformanceMe{
		service: service,
	}
}

func (h *GetPerformanceMe) GetPerformanceMe(c *gin.Context) {
	// Parse query parameters
	var req adminReportModel.GetPerformanceMeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Parse and validate dates
	startDate, err := utils.DateStringToTime(req.StartDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"startDate": "startDate 轉換類型失敗",
		})
		return
	}

	endDate, err := utils.DateStringToTime(req.EndDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"endDate": "endDate 轉換類型失敗",
		})
		return
	}

	parsedReq := adminReportModel.GetPerformanceMeParsedRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Get staff context from JWT middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Extract store IDs from staff context
	storeIds := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeIds[i] = store.ID
	}

	// Call service
	performance, err := h.service.GetPerformanceMe(c.Request.Context(), parsedReq, staffContext.UserID, storeIds)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(performance))
}
