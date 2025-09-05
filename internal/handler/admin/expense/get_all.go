package adminExpense

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminExpenseService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminExpenseService.GetAllInterface
}

func NewGetAll(service adminExpenseService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	storeIDStr := c.Param("storeId")
	if storeIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId 為必填項目",
		})
		return
	}
	storeID, err := utils.ParseID(storeIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	var req adminExpenseModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim category
	if req.Category != nil {
		*req.Category = strings.TrimSpace(*req.Category)
	}

	var supplierID *int64
	var payerID *int64
	if req.SupplierID != nil {
		parsed, err := utils.ParseID(*req.SupplierID)
		supplierID = &parsed
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"supplierId": "supplierId 類型轉換失敗",
			})
			return
		}
	}
	if req.PayerID != nil {
		parsed, err := utils.ParseID(*req.PayerID)
		payerID = &parsed
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"payerId": "payerId 類型轉換失敗",
			})
			return
		}
	}

	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	staff, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}
	storeIDs := make([]int64, len(staff.StoreList))
	for i, store := range staff.StoreList {
		storeIDs[i] = store.ID
	}

	parsedReq := adminExpenseModel.GetAllParsedRequest{
		Category:     req.Category,
		SupplierID:   supplierID,
		PayerID:      payerID,
		IsReimbursed: req.IsReimbursed,
		Limit:        limit,
		Offset:       offset,
		Sort:         sort,
	}

	response, err := h.service.GetAll(c.Request.Context(), storeID, parsedReq, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
