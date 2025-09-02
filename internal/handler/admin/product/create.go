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

type Create struct {
	service adminProductService.CreateInterface
}

func NewCreate(service adminProductService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
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

	var req adminProductModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name, unit, storageLocation, note
	req.Name = strings.TrimSpace(req.Name)
	if req.Unit != nil {
		trimmed := strings.TrimSpace(*req.Unit)
		req.Unit = &trimmed
	}
	if req.StorageLocation != nil {
		trimmed := strings.TrimSpace(*req.StorageLocation)
		req.StorageLocation = &trimmed
	}
	if req.Note != nil {
		trimmed := strings.TrimSpace(*req.Note)
		req.Note = &trimmed
	}

	brandID, err := utils.ParseID(req.BrandID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"brandId": "brandId 類型轉換失敗",
		})
		return
	}
	categoryID, err := utils.ParseID(req.CategoryID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"categoryId": "categoryId 類型轉換失敗",
		})
		return
	}

	parsedCurrentStock := int(*req.CurrentStock)
	parsedSafetyStock := -1
	if req.SafetyStock != nil {
		parsedSafetyStock = int(*req.SafetyStock)
	}

	parsedReq := adminProductModel.CreateParsedRequest{
		Name:            req.Name,
		BrandID:         brandID,
		CategoryID:      categoryID,
		CurrentStock:    parsedCurrentStock,
		SafetyStock:     parsedSafetyStock,
		Unit:            req.Unit,
		StorageLocation: req.StorageLocation,
		Note:            req.Note,
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	response, err := h.service.Create(c.Request.Context(), storeID, parsedReq, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
