package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service bookingService.CreateInterface
}

func NewCreate(service bookingService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

// Create handles POST /api/bookings/me
func (h *Create) Create(c *gin.Context) {
	var req bookingModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	storeId, err := utils.ParseID(req.StoreId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	stylistId, err := utils.ParseID(req.StylistId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"stylistId": "stylistId 類型轉換失敗",
		})
		return
	}

	timeSlotId, err := utils.ParseID(req.TimeSlotId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"timeSlotId": "timeSlotId 類型轉換失敗",
		})
		return
	}

	mainServiceId, err := utils.ParseID(req.MainServiceId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"mainServiceId": "mainServiceId 類型轉換失敗",
		})
		return
	}

	subServiceIds := make([]int64, len(req.SubServiceIds))
	for i, subServiceId := range req.SubServiceIds {
		subServiceId, err := utils.ParseID(subServiceId)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"subServiceId": "subServiceId 類型轉換失敗",
			})
			return
		}
		subServiceIds[i] = subServiceId
	}

	parsedRequest := bookingModel.CreateParsedRequest{
		StoreId:       storeId,
		StylistId:     stylistId,
		TimeSlotId:    timeSlotId,
		MainServiceId: mainServiceId,
		SubServiceIds: subServiceIds,
		Note:          req.Note,
		IsChatEnabled: req.IsChatEnabled,
	}

	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.service.Create(c.Request.Context(), parsedRequest, customerContext.CustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
