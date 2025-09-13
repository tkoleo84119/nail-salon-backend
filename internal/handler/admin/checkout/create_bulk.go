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

type CreateBulk struct {
	service adminCheckoutService.CreateBulkInterface
}

func NewCreateBulk(service adminCheckoutService.CreateBulkInterface) *CreateBulk {
	return &CreateBulk{
		service: service,
	}
}

func (h *CreateBulk) CreateBulk(c *gin.Context) {
	storeID := c.Param("storeId")
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

	var req adminCheckoutModel.CreateBulkRequest
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

	checkouts := make([]adminCheckoutModel.CreateBulkParsedCheckoutItems, len(req.Checkouts))
	for i, checkout := range req.Checkouts {
		parsedBookingID, err := utils.ParseID(checkout.BookingID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"bookingID": "bookingID 類型轉換失敗",
			})
			return
		}

		details := make([]adminCheckoutModel.CreateBulkParsedDetailItems, len(checkout.Details))
		applyCount := 0
		for j, detail := range checkout.Details {
			parsedDetailID, err := utils.ParseID(detail.ID)
			if err != nil {
				errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
					"detailID": "detailID 類型轉換失敗",
				})
			}

			useCoupon := false
			if detail.UseCoupon != nil {
				useCoupon = *detail.UseCoupon
				if useCoupon {
					applyCount++
				}
			}

			details[j] = adminCheckoutModel.CreateBulkParsedDetailItems{
				ID:        parsedDetailID,
				Price:     detail.Price,
				UseCoupon: useCoupon,
			}
		}

		checkouts[i] = adminCheckoutModel.CreateBulkParsedCheckoutItems{
			BookingID:  parsedBookingID,
			PaidAmount: checkout.PaidAmount,
			ApplyCount: int64(applyCount),
			Details:    details,
		}
	}

	parsedRequest := adminCheckoutModel.CreateBulkParsedRequest{
		PaymentMethod:    req.PaymentMethod,
		CustomerCouponID: customerCouponID,
		Checkouts:        checkouts,
	}

	// Get staff context from JWT middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.service.CreateBulk(c.Request.Context(), parsedStoreID, parsedRequest, staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
