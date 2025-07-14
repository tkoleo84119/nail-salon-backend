package staff

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

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
			c.JSON(http.StatusBadRequest, common.ValidationErrorResponse(validationErrors))
		} else {
			errors := map[string]string{"request": "JSON格式錯誤"}
			c.JSON(http.StatusBadRequest, common.ErrorResponse("請求錯誤", errors))
		}
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errors := map[string]string{"token": "無法取得使用者認證資訊"}
		c.JSON(http.StatusUnauthorized, common.ErrorResponse("認證失敗", errors))
		return
	}

	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	response, err := h.createStaffService.CreateStaff(c.Request.Context(), req, staffContext.Role, creatorStoreIDs)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}

// handleServiceError 處理service層的錯誤並轉換為適當的中文回應
func (h *CreateStaffHandler) handleServiceError(c *gin.Context, err error) {
	errMsg := err.Error()

	if h.isValidationError(errMsg) {
		errors := h.getValidationErrorMessage(errMsg)
		c.JSON(http.StatusBadRequest, common.ErrorResponse("輸入驗證失敗", errors))
		return
	}

	if h.isPermissionError(errMsg) {
		errors := map[string]string{"permission": h.getPermissionErrorMessage(errMsg)}
		c.JSON(http.StatusForbidden, common.ErrorResponse("權限不足", errors))
		return
	}

	errors := map[string]string{"server": "建立帳號時發生錯誤"}
	c.JSON(http.StatusInternalServerError, common.ErrorResponse("系統錯誤", errors))
}

func (h *CreateStaffHandler) isValidationError(errMsg string) bool {
	validationErrors := []string{
		"username or email already exists",
		"some stores do not exist",
		"some stores are not active",
		"invalid role:",
		"cannot create SUPER_ADMIN role",
	}

	for _, validationErr := range validationErrors {
		if strings.Contains(errMsg, validationErr) {
			return true
		}
	}
	return false
}

func (h *CreateStaffHandler) isPermissionError(errMsg string) bool {
	permissionErrors := []string{
		"insufficient permissions to create staff",
		"SUPER_ADMIN cannot create another SUPER_ADMIN",
		"ADMIN can only create MANAGER or STYLIST roles",
		"no permission to assign store",
	}

	for _, permissionErr := range permissionErrors {
		if strings.Contains(errMsg, permissionErr) {
			return true
		}
	}
	return false
}

func (h *CreateStaffHandler) getValidationErrorMessage(errMsg string) map[string]string {
	errors := make(map[string]string)

	switch {
	case strings.Contains(errMsg, "username or email already exists"):
		errors["username"] = "帳號已存在"
		errors["email"] = "Email已存在"
	case strings.Contains(errMsg, "some stores do not exist"):
		errors["store_ids"] = "部分門市不存在"
	case strings.Contains(errMsg, "some stores are not active"):
		errors["store_ids"] = "部分門市未啟用"
	case strings.Contains(errMsg, "invalid role"):
		errors["role"] = "無效的角色"
	case strings.Contains(errMsg, "cannot create SUPER_ADMIN role"):
		errors["role"] = "不可建立超級管理員帳號"
	default:
		errors["validation"] = "輸入資料有誤"
	}

	return errors
}

func (h *CreateStaffHandler) getPermissionErrorMessage(errMsg string) string {
	switch {
	case strings.Contains(errMsg, "insufficient permissions to create staff"):
		return "權限不足，無法建立員工帳號"
	case strings.Contains(errMsg, "SUPER_ADMIN cannot create another SUPER_ADMIN"):
		return "超級管理員無法建立另一個超級管理員"
	case strings.Contains(errMsg, "ADMIN can only create MANAGER or STYLIST roles"):
		return "管理員只能建立主管或美甲師帳號"
	case strings.Contains(errMsg, "no permission to assign store"):
		return "沒有權限指派此門市"
	default:
		return "權限不足"
	}
}