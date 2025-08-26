package adminCheckout

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminCheckoutModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/checkout"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCheckoutService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/checkout"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminCheckoutService.CreateInterface
}

func NewCreate(service adminCheckoutService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	storeID := c.Param("storeID")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeID": "storeID 為必填項目",
		})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeID": "storeID 類型轉換失敗",
		})
		return
	}

	bookingID := c.Param("bookingID")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"bookingID": "bookingID 為必填項目",
		})
		return
	}
	parsedBookingID, err := utils.ParseID(bookingID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"bookingID": "bookingID 類型轉換失敗",
		})
		return
	}

	var req adminCheckoutModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	var customerCouponID *int64
	if req.CustomerCouponID != nil {
		parsedCustomerCouponID, err := utils.ParseID(*req.CustomerCouponID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"customerCouponID": "customerCouponID 類型轉換失敗",
			})
			return
		}
		customerCouponID = &parsedCustomerCouponID
	}

	bookingDetails := make([]adminCheckoutModel.CreateBookingDetailParsed, len(req.BookingDetails))
	for i, bookingDetail := range req.BookingDetails {
		parsedBookingDetailID, err := utils.ParseID(bookingDetail.ID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"bookingDetailID": "bookingDetailID 類型轉換失敗",
			})
			return
		}
		bookingDetails[i] = adminCheckoutModel.CreateBookingDetailParsed{
			ID:        parsedBookingDetailID,
			Price:     bookingDetail.Price,
			UseCoupon: bookingDetail.UseCoupon,
		}
	}

	parsedRequest := adminCheckoutModel.CreateParsedRequest{
		PaymentMethod:    req.PaymentMethod,
		CustomerCouponID: customerCouponID,
		PaidAmount:       req.PaidAmount,
		BookingDetails:   bookingDetails,
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

	response, err := h.service.Create(c.Request.Context(), parsedStoreID, parsedBookingID, parsedRequest, staffContext.UserID, storeIds)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
