package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	serviceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	serviceService "github.com/tkoleo84119/nail-salon-backend/internal/service/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service serviceService.GetAllInterface
}

func NewGetAll(service serviceService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Path parameter validation
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

	// Query parameter validation
	var queryParams serviceModel.GetAllRequest
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	limit, offset := utils.SetDefaultValuesOfPagination(queryParams.Limit, queryParams.Offset, 20, 0)
	sort := utils.TransformSort(queryParams.Sort)

	parsedQueryParams := serviceModel.GetAllParsedRequest{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	// Service layer call (no authentication required per spec)
	response, err := h.service.GetAll(c.Request.Context(), parsedStoreID, parsedQueryParams)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
