package adminExpense

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminExpenseService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminExpenseService.CreateInterface
}

func NewCreate(service adminExpenseService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
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

	var req adminExpenseModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	req.Category = strings.TrimSpace(req.Category)
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	expenseDate, err := utils.DateStringToTime(req.ExpenseDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
			"expenseDate": "expenseDate 日期格式錯誤，應為 YYYY-MM-DD",
		})
		return
	}

	var parsedPayerID *int64
	if req.PayerID != nil && *req.PayerID != "" {
		payerID, err := utils.ParseID(*req.PayerID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"payerId": "payerId 轉換類型失敗",
			})
			return
		}
		parsedPayerID = &payerID
	}

	var parsedSupplierID *int64
	if req.SupplierID != nil && *req.SupplierID != "" {
		supplierID, err := utils.ParseID(*req.SupplierID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"supplierId": "supplierId 轉換類型失敗",
			})
			return
		}
		parsedSupplierID = &supplierID
	}

	amount := int64(0)
	if req.Amount != nil {
		amount = *req.Amount
	}

	items := make([]adminExpenseModel.CreateItemParsedRequest, 0)
	for _, item := range req.Items {
		productID, err := utils.ParseID(item.ProductID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"productId": "productId 轉換類型失敗",
			})
			return
		}

		quantity := int64(0)
		if item.Quantity != nil {
			quantity = *item.Quantity
		}

		price := int64(0)
		if item.Price != nil {
			price = *item.Price
		}

		expirationDate := time.Time{}
		if item.ExpirationDate != nil {
			expirationDate, err = utils.DateStringToTime(*item.ExpirationDate)
			if err != nil {
				errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
					"expirationDate": "expirationDate 日期格式錯誤，應為 YYYY-MM-DD",
				})
				return
			}
		}

		isArrived := false
		if item.IsArrived != nil {
			isArrived = *item.IsArrived
		}

		arrivalDate := time.Time{}
		if item.ArrivalDate != nil {
			arrivalDate, err = utils.DateStringToTime(*item.ArrivalDate)
			if err != nil {
				errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
					"arrivalDate": "arrivalDate 日期格式錯誤，應為 YYYY-MM-DD",
				})
				return
			}
		}

		if item.StorageLocation != nil {
			trimmed := strings.TrimSpace(*item.StorageLocation)
			item.StorageLocation = &trimmed
		}

		if item.Note != nil {
			trimmed := strings.TrimSpace(*item.Note)
			item.Note = &trimmed
		}

		items = append(items, adminExpenseModel.CreateItemParsedRequest{
			ProductID:       productID,
			Quantity:        quantity,
			Price:           price,
			ExpirationDate:  &expirationDate,
			IsArrived:       isArrived,
			ArrivalDate:     &arrivalDate,
			StorageLocation: item.StorageLocation,
			Note:            item.Note,
		})
	}

	parsedReq := adminExpenseModel.CreateParsedRequest{
		SupplierID:  parsedSupplierID,
		Category:    req.Category,
		Amount:      amount,
		OtherFee:    req.OtherFee,
		ExpenseDate: expenseDate,
		Note:        req.Note,
		PayerID:     parsedPayerID,
		Items:       items,
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

	// Call service layer
	response, err := h.service.Create(c.Request.Context(), parsedStoreID, parsedReq, staffContext.UserID, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
