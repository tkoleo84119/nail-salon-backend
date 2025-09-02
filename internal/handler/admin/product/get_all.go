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

type GetAll struct {
	service adminProductService.GetAllInterface
}

func NewGetAll(service adminProductService.GetAllInterface) *GetAll {
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

	var req adminProductModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	var brandID int64
	var categoryID int64
	if req.BrandID != nil && *req.BrandID != "" {
		brandID, err = utils.ParseID(*req.BrandID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"brandId": "brandId 類型轉換失敗",
			})
			return
		}
	}
	if req.CategoryID != nil && *req.CategoryID != "" {
		categoryID, err = utils.ParseID(*req.CategoryID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"categoryId": "categoryId 類型轉換失敗",
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

	parsedReq := adminProductModel.GetAllParsedRequest{
		BrandID:             &brandID,
		CategoryID:          &categoryID,
		Name:                req.Name,
		LessThanSafetyStock: req.LessThanSafetyStock,
		Limit:               limit,
		Offset:              offset,
		Sort:                sort,
	}

	response, err := h.service.GetAll(c.Request.Context(), storeID, parsedReq, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
