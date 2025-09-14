package adminBooking

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminBookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminBookingService.GetAllInterface
}

func NewGetAll(service adminBookingService.GetAllInterface) *GetAll {
	return &GetAll{service: service}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Get path parameter
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"storeId": "storeId 為必填項目"})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"storeId": "storeId 轉換類型失敗"})
		return
	}

	// Parse query parameters
	var req adminBookingModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// set limit and offset
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	var stylistID *int64
	if req.StylistID != nil {
		parsedStylistID, err := utils.ParseID(*req.StylistID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"stylistId": "stylistId 轉換類型失敗"})
			return
		}
		stylistID = &parsedStylistID
	}

	var customerID *int64
	if req.CustomerID != nil {
		parsedCustomerID, err := utils.ParseID(*req.CustomerID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"customerId": "customerId 轉換類型失敗"})
			return
		}
		customerID = &parsedCustomerID
	}

	var startDate *time.Time
	if req.StartDate != nil {
		parsedStartDate, err := utils.DateStringToTime(*req.StartDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"startDate": "startDate 轉換類型失敗"})
			return
		}
		startDate = &parsedStartDate
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := utils.DateStringToTime(*req.EndDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"endDate": "endDate 轉換類型失敗"})
			return
		}
		endDate = &parsedEndDate
	}

	parsedReq := adminBookingModel.GetAllParsedRequest{
		StylistID:  stylistID,
		CustomerID: customerID,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     req.Status,
		Limit:      limit,
		Offset:     offset,
		Sort:       sort,
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
	bookings, err := h.service.GetAll(c.Request.Context(), parsedStoreID, parsedReq, staffContext.Role, storeIds)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(bookings))
}
