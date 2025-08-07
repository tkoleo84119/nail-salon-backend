package adminStylist

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStylistListHandler struct {
	service adminStylistService.GetStylistListServiceInterface
}

func NewGetStylistListHandler(service adminStylistService.GetStylistListServiceInterface) *GetStylistListHandler {
	return &GetStylistListHandler{
		service: service,
	}
}

func (h *GetStylistListHandler) GetStylistList(c *gin.Context) {
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
	var req adminStylistModel.GetStylistListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// set default limit and offset
	limit := 20
	offset := 0
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	// parse sort
	sort := []string{}
	if req.Sort != nil && *req.Sort != "" {
		sort = strings.Split(*req.Sort, ",")
	}

	parsedReq := adminStylistModel.GetStylistListParsedRequest{
		Name:        req.Name,
		IsIntrovert: req.IsIntrovert,
		Limit:       limit,
		Offset:      offset,
		Sort:        sort,
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
	response, err := h.service.GetStylistList(c.Request.Context(), parsedStoreID, parsedReq, staffContext.Role, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
