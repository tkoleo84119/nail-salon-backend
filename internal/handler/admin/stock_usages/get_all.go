package adminStockUsages

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStockUsagesService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminStockUsagesService.GetAllInterface
}

func NewGetAll(service adminStockUsagesService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Get path parameter
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
			"storeId": "storeId 轉換類型失敗",
		})
		return
	}

	// Parse and validate query parameters
	var req adminStockUsagesModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	// set default value
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	parsedReq := adminStockUsagesModel.GetAllParsedRequest{
		Name:    req.Name,
		IsInUse: req.IsInUse,
		Limit:   limit,
		Offset:  offset,
		Sort:    sort,
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

	// Service layer call
	response, err := h.service.GetAll(c.Request.Context(), parsedStoreID, parsedReq, staffContext.Role, storeIds)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
