package adminAccount

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminAccountModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAccountService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/account"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminAccountService.CreateInterface
}

func NewCreate(service adminAccountService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	// Parse and validate request
	var req adminAccountModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	parsedStoreID, err := utils.ParseID(req.StoreID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	// trim name, note
	req.Name = strings.TrimSpace(req.Name)
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	parsedReq := adminAccountModel.CreateParsedRequest{
		StoreID: parsedStoreID,
		Name:    req.Name,
		Note:    req.Note,
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

	// Call service
	response, err := h.service.Create(c.Request.Context(), parsedReq, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
