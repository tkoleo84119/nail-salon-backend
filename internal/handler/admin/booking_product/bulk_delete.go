package adminBookingProduct

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminBookingProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking_product"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminBookingProductService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking_product"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type BulkDelete struct {
	service adminBookingProductService.BulkDeleteInterface
}

func NewBulkDelete(service adminBookingProductService.BulkDeleteInterface) *BulkDelete {
	return &BulkDelete{
		service: service,
	}
}

func (h *BulkDelete) BulkDelete(c *gin.Context) {
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

	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"bookingId": "bookingId 為必填項目",
		})
		return
	}
	parsedBookingID, err := utils.ParseID(bookingID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"bookingId": "bookingId 類型轉換失敗",
		})
		return
	}

	var req adminBookingProductModel.BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	productIds := make([]int64, len(req.ProductIds))
	for i, productID := range req.ProductIds {
		parsedProductID, err := utils.ParseID(productID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"productIds": "productIds 類型轉換失敗",
			})
			return
		}
		productIds[i] = parsedProductID
	}

	parsedRequest := adminBookingProductModel.BulkDeleteParsedRequest{
		ProductIds: productIds,
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	staffStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		staffStoreIDs[i] = store.ID
	}

	response, err := h.service.BulkDelete(c.Request.Context(), parsedStoreID, parsedBookingID, parsedRequest, staffContext.Role, staffStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
