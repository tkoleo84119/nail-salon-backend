package adminProduct

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminProductService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/product"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminProductService.UpdateInterface
}

func NewUpdate(service adminProductService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
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

	productIDStr := c.Param("productId")
	if productIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"productId": "productId 為必填項目",
		})
		return
	}
	productID, err := utils.ParseID(productIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"productId": "productId 類型轉換失敗",
		})
		return
	}

	var req adminProductModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// trim name
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	var brandID *int64
	var categoryID *int64
	if req.BrandID != nil && *req.BrandID != "" {
		parsed, err := utils.ParseID(*req.BrandID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"brandId": "brandId 類型轉換失敗",
			})
			return
		}
		brandID = &parsed
	}
	if req.CategoryID != nil && *req.CategoryID != "" {
		parsed, err := utils.ParseID(*req.CategoryID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"categoryId": "categoryId 類型轉換失敗",
			})
			return
		}
		categoryID = &parsed
	}

	var currentStock *int
	var safetyStock *int
	if req.CurrentStock != nil {
		converted := int(*req.CurrentStock)
		currentStock = &converted
	}
	if req.SafetyStock != nil {
		converted := int(*req.SafetyStock)
		safetyStock = &converted
	}

	parsedReq := adminProductModel.UpdateParsedRequest{
		BrandID:         brandID,
		CategoryID:      categoryID,
		Name:            req.Name,
		CurrentStock:    currentStock,
		SafetyStock:     safetyStock,
		Unit:            req.Unit,
		StorageLocation: req.StorageLocation,
		Note:            req.Note,
		IsActive:        req.IsActive,
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	response, err := h.service.Update(c.Request.Context(), storeID, productID, parsedReq, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
