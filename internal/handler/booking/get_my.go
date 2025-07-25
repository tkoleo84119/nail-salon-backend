package booking

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMyBookingsHandler struct {
	service bookingService.GetMyBookingsServiceInterface
}

func NewGetMyBookingsHandler(service bookingService.GetMyBookingsServiceInterface) *GetMyBookingsHandler {
	return &GetMyBookingsHandler{
		service: service,
	}
}

func (h *GetMyBookingsHandler) GetMyBookings(c *gin.Context) {
	// Query parameter validation
	var queryParams bookingModel.GetMyBookingsQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Parse status parameter if provided (supports comma-separated values)
	statusParam := c.Query("status")
	if statusParam != "" {
		// Split comma-separated statuses and trim spaces
		statuses := strings.Split(statusParam, ",")
		for i, status := range statuses {
			statuses[i] = strings.TrimSpace(status)
		}
		queryParams.Status = statuses
	}

	// Set defaults if not provided
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}

	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.GetMyBookings(c.Request.Context(), queryParams, *customerContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}