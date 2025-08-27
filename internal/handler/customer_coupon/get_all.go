package customerCoupon

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	customerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer_coupon"
	customerCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service customerCouponService.GetAllInterface
}

func NewGetAll(service customerCouponService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	var req customerCouponModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	parsedReq := customerCouponModel.GetAllParsedRequest{
		IsUsed: req.IsUsed,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	resp, err := h.service.GetAll(c.Request.Context(), customerContext.CustomerID, parsedReq)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(resp))
}
