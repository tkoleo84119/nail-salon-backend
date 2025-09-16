package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type LineVerifyResponse struct {
	Iss     string   `json:"iss"`
	Sub     string   `json:"sub"`
	Aud     string   `json:"aud"`
	Exp     int64    `json:"exp"`
	Iat     int64    `json:"iat"`
	Nonce   string   `json:"nonce,omitempty"`
	Amr     []string `json:"amr,omitempty"`
	Name    string   `json:"name"`
	Picture string   `json:"picture,omitempty"`
	Email   string   `json:"email,omitempty"`
}

type LineValidator struct {
	channelID      string
	verifyEndpoint string
	httpClient     *http.Client
}

type LineMessageClient struct {
	channelAccessToken string
	httpClient         *http.Client
	messageEndpoint    string
}

type TextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type FlexMessage struct {
	Type     string      `json:"type"`
	AltText  string      `json:"altText"`
	Contents interface{} `json:"contents"`
}

type PushMessageRequest struct {
	To       string        `json:"to"`
	Messages []interface{} `json:"messages"`
}

type BookingData struct {
	StoreName       string   `json:"storeName"`
	StoreAddress    string   `json:"storeAddress"`
	Date            string   `json:"date"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	CustomerName    *string  `json:"customerName,omitempty"`
	CustomerPhone   *string  `json:"customerPhone,omitempty"`
	StylistName     string   `json:"stylistName"`
	MainServiceName string   `json:"mainServiceName"`
	SubServiceNames []string `json:"subServiceNames,omitempty"`
}

func NewLineValidator(channelID string) *LineValidator {
	return &LineValidator{
		channelID:      channelID,
		verifyEndpoint: "https://api.line.me/oauth2/v2.1/verify",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (v *LineValidator) ValidateIdToken(idToken string) (*common.LineProfile, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, errorCodes.NewServiceError(errorCodes.AuthLineTokenInvalid, "invalid token structure", nil)
	}

	data := url.Values{}
	data.Set("id_token", idToken)
	data.Set("client_id", v.channelID)

	req, err := http.NewRequest("POST", v.verifyEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to verify token", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to read response", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if errorDesc, ok := errorResp["error_description"].(string); ok {
				if strings.Contains(strings.ToLower(errorDesc), "expired") {
					return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthLineTokenExpired)
				}
			}
		}
		return nil, errorCodes.NewServiceError(
			errorCodes.AuthLineTokenInvalid,
			fmt.Sprintf("verification failed with status %d: %s", resp.StatusCode, string(body)),
			nil,
		)
	}

	var verifyResp LineVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to parse verification response", err)
	}

	if verifyResp.Iss != "https://access.line.me" {
		return nil, errorCodes.NewServiceError(errorCodes.AuthLineTokenInvalid, "invalid issuer", nil)
	}

	if verifyResp.Aud != v.channelID {
		return nil, errorCodes.NewServiceError(errorCodes.AuthLineTokenInvalid, "invalid audience", nil)
	}

	profile := &common.LineProfile{
		ProviderUid: verifyResp.Sub,
		Name:        verifyResp.Name,
	}

	if verifyResp.Email != "" {
		profile.Email = &verifyResp.Email
	}

	return profile, nil
}

func NewLineMessenger(channelAccessToken string) *LineMessageClient {
	return &LineMessageClient{
		channelAccessToken: channelAccessToken,
		messageEndpoint:    "https://api.line.me/v2/bot/message/push",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendTextMessage sends a text message to a user
func (c *LineMessageClient) SendTextMessage(userID string, text string) error {
	message := TextMessage{
		Type: "text",
		Text: text,
	}

	return c.sendMessage(userID, message)
}

// SendFlexMessage sends a flex message to a user, detail contents can be found in https://developers.line.biz/en/docs/messaging-api/using-flex-messages/
func (c *LineMessageClient) SendFlexMessage(userID string, altText string, contents interface{}) error {
	message := FlexMessage{
		Type:     "flex",
		AltText:  altText,
		Contents: contents,
	}

	return c.sendMessage(userID, message)
}

// SendBookingNotification sends a booking notification to a user
func (c *LineMessageClient) SendBookingNotification(userID string, action common.BookingAction, bookingData *BookingData) error {
	actionText := c.getActionText(action)
	flexContent := c.buildBookingFlexContent(bookingData, action, actionText)

	altText := fmt.Sprintf("%s - %s", bookingData.StoreName, actionText)
	return c.SendFlexMessage(userID, altText, flexContent)
}

// sendMessage is basic function to send message to line
func (c *LineMessageClient) sendMessage(userID string, message interface{}) error {
	requestData := PushMessageRequest{
		To:       userID,
		Messages: []interface{}{message},
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to marshal message data", err)
	}

	req, err := http.NewRequest("POST", c.messageEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.channelAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to send message", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errorCodes.NewServiceError(
			errorCodes.SysInternalError,
			fmt.Sprintf("LINE API error: status %d, body: %s", resp.StatusCode, string(body)),
			nil,
		)
	}

	return nil
}

// getActionText is function to get action text
func (c *LineMessageClient) getActionText(action common.BookingAction) string {
	switch action {
	case common.BookingActionCreated:
		return "預約成功"
	case common.BookingActionUpdated:
		return "預約更改"
	case common.BookingActionCancelled:
		return "預約取消"
	default:
		return "預約通知"
	}
}

// buildBookingFlexContent is function to build booking flex content
func (c *LineMessageClient) buildBookingFlexContent(bookingData *BookingData, _ common.BookingAction, actionText string) map[string]interface{} {
	customerName := "顧客"
	if bookingData.CustomerName != nil && *bookingData.CustomerName != "" {
		customerName = *bookingData.CustomerName
	}

	customerPhone := "未提供"
	if bookingData.CustomerPhone != nil && *bookingData.CustomerPhone != "" {
		customerPhone = *bookingData.CustomerPhone
	}

	subServicesText := "無"
	if len(bookingData.SubServiceNames) > 0 {
		subServicesText = strings.Join(bookingData.SubServiceNames, ", ")
	}

	header := map[string]interface{}{
		"type":   "box",
		"layout": "vertical",
		"contents": []map[string]interface{}{
			{
				"type":   "text",
				"text":   bookingData.StoreName,
				"weight": "bold",
				"size":   "xl",
				"color":  "#1DB446",
			},
			{
				"type":   "text",
				"text":   actionText,
				"size":   "md",
				"weight": "bold",
			},
			{
				"type":  "text",
				"text":  formatDateTimeWithWeekday(bookingData.Date, bookingData.StartTime, bookingData.EndTime),
				"size":  "xs",
				"color": "#666666",
			},
		},
		"spacing": "sm",
	}

	bodyContents := []map[string]interface{}{
		{
			"type":   "box",
			"layout": "vertical",
			"margin": "none",
			"contents": []map[string]interface{}{
				{
					"type":   "box",
					"layout": "horizontal",
					"contents": []map[string]interface{}{
						{
							"type":  "text",
							"text":  "預約姓名",
							"size":  "sm",
							"color": "#666666",
							"flex":  0,
						},
						{
							"type":  "text",
							"text":  customerName,
							"size":  "sm",
							"align": "end",
						},
					},
				},
				{
					"type":   "box",
					"layout": "horizontal",
					"contents": []map[string]interface{}{
						{
							"type":  "text",
							"text":  "預約者電話",
							"size":  "sm",
							"color": "#666666",
							"flex":  0,
						},
						{
							"type":  "text",
							"text":  customerPhone,
							"size":  "sm",
							"align": "end",
						},
					},
				},
				{
					"type":   "box",
					"layout": "horizontal",
					"contents": []map[string]interface{}{
						{
							"type":  "text",
							"text":  "設計師",
							"size":  "sm",
							"color": "#666666",
							"flex":  0,
						},
						{
							"type":  "text",
							"text":  bookingData.StylistName,
							"size":  "sm",
							"align": "end",
						},
					},
				},
			},
			"spacing": "sm",
		},
		{
			"type":   "separator",
			"margin": "xxl",
		},
		{
			"type":   "box",
			"layout": "vertical",
			"contents": []map[string]interface{}{
				{
					"type":   "box",
					"layout": "horizontal",
					"contents": []map[string]interface{}{
						{
							"type":  "text",
							"text":  "服務項目",
							"size":  "sm",
							"color": "#666666",
							"flex":  0,
						},
						{
							"type":  "text",
							"text":  bookingData.MainServiceName,
							"size":  "sm",
							"align": "end",
						},
					},
				},
				{
					"type":   "box",
					"layout": "vertical",
					"contents": []map[string]interface{}{
						{
							"type":  "text",
							"text":  subServicesText,
							"size":  "sm",
							"color": "#666666",
							"wrap":  true,
							"align": "end",
						},
					},
				},
			},
			"margin": "xxl",
		},
	}

	body := map[string]interface{}{
		"type":     "box",
		"layout":   "vertical",
		"spacing":  "md",
		"contents": bodyContents,
	}

	return map[string]interface{}{
		"type":   "bubble",
		"header": header,
		"body":   body,
		"styles": map[string]interface{}{
			"header": map[string]interface{}{
				"separator": false,
			},
			"body": map[string]interface{}{
				"separator": true,
			},
			"footer": map[string]interface{}{
				"separator": false,
			},
		},
	}
}
