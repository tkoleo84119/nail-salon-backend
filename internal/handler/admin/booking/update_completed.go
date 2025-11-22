package adminBooking

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminBookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateCompleted struct {
	service adminBookingService.UpdateCompletedInterface
}

func NewUpdateCompleted(service adminBookingService.UpdateCompletedInterface) *UpdateCompleted {
	return &UpdateCompleted{
		service: service,
	}
}

func (h *UpdateCompleted) UpdateCompleted(c *gin.Context) {
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
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	bookingID := c.Param("bookingId")
	if bookingID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"bookingId": "bookingId 為必填項目",
		})
		return
	}
	parsedBookingID, err := utils.ParseID(bookingID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"bookingId": "bookingId 類型轉換失敗",
		})
		return
	}

	// Bind request body
	var req adminBookingModel.UpdateCompletedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// trim url and validate url format
	if req.PinterestImageUrls != nil && len(*req.PinterestImageUrls) > 0 {
		rawUrls := make([]string, len(*req.PinterestImageUrls))

		type result struct {
			index int
			url   string
			err   error
		}
		resultChan := make(chan result, len(*req.PinterestImageUrls))

		for i, url := range *req.PinterestImageUrls {
			go func(idx int, u string) {
				trimmedUrl := strings.TrimSpace(u)

				// if url is already a valid url (ex: https://www.pinterest.com/pin/xxxxx/), return it
				if strings.HasPrefix(trimmedUrl, "https://www.pinterest.com/pin/") {
					resultChan <- result{
						index: idx,
						url:   trimmedUrl,
						err:   nil,
					}
					return
				}

				// valid url format => https://pin.it/xxxxx
				if !strings.HasPrefix(trimmedUrl, "https://pin.it/") {
					resultChan <- result{index: idx, err: fmt.Errorf("格式錯誤")}
					return
				}
				rawUrl, err := utils.ResolveURL(trimmedUrl)
				if err != nil {
					resultChan <- result{index: idx, err: err}
					return
				}

				// validate raw url format => https://www.pinterest.com/pin/xxxxx/
				if !strings.HasPrefix(rawUrl, "https://www.pinterest.com/pin/") {
					resultChan <- result{index: idx, err: fmt.Errorf("無法解析")}
					return
				}

				resultChan <- result{index: idx, url: rawUrl}
			}(i, url)
		}

		// get results from channel
		for i := 0; i < len(*req.PinterestImageUrls); i++ {
			res := <-resultChan
			if res.err != nil {
				errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
					"pinterestImageUrls": "無法解析 Pinterest 圖片 URL",
				})
				return
			}
			rawUrls[res.index] = res.url
		}

		req.PinterestImageUrls = &rawUrls
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	updaterStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		updaterStoreIDs[i] = store.ID
	}

	// Call service
	response, err := h.service.UpdateCompleted(c.Request.Context(), parsedStoreID, parsedBookingID, req, staffContext.Role, updaterStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
