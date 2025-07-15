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
	createStaffService staffService.CreateStaffServiceInterface
}

func NewCreateStaffHandler(createStaffService staffService.CreateStaffServiceInterface) *CreateStaffHandler {
	return &CreateStaffHandler{
		createStaffService: createStaffService,
	}
}

func (h *CreateStaffHandler) CreateStaff(c *gin.Context) {
	var req staff.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsValidationError(err) {
			validationErrors := utils.ExtractValidationErrors(err)
			errorCodes.RespondWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		} else {
			fieldErrors := map[string]string{"request": "JSON格式錯誤"}
			errorCodes.RespondWithError(c, errorCodes.ValJsonFormat, fieldErrors)
		}
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.RespondWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	response, err := h.createStaffService.CreateStaff(c.Request.Context(), req, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}

