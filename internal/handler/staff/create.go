package staff

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStaffHandler struct {
	service staffService.CreateStaffServiceInterface
}

func NewCreateStaffHandler(service staffService.CreateStaffServiceInterface) *CreateStaffHandler {
	return &CreateStaffHandler{
		service: service,
	}
}

func (h *CreateStaffHandler) CreateStaff(c *gin.Context) {
	var req staff.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.RespondWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeID, err := utils.ParseID(store.ID)
		if err != nil {
			errorCodes.RespondWithError(c, errorCodes.AuthContextMissing, nil)
			return
		}
		creatorStoreIDs[i] = storeID
	}

	response, err := h.service.CreateStaff(c.Request.Context(), req, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
