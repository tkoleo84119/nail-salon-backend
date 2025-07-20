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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
)

type LineVerifyResponse struct {
	Iss     string `json:"iss"`
	Sub     string `json:"sub"`
	Aud     string `json:"aud"`
	Exp     int64  `json:"exp"`
	Iat     int64  `json:"iat"`
	Nonce   string `json:"nonce,omitempty"`
	Amr     []string `json:"amr,omitempty"`
	Name    string `json:"name"`
	Picture string `json:"picture,omitempty"`
	Email   string `json:"email,omitempty"`
}

type LineValidator struct {
	channelID      string
	verifyEndpoint string
	httpClient     *http.Client
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

func (v *LineValidator) ValidateIdToken(idToken string) (*customer.CustomerProfile, error) {
	if v.channelID == "YOUR_LINE_CHANNEL_ID" || v.channelID == "" {
		return MockValidateLineIdToken(idToken)
	}

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

	profile := &customer.CustomerProfile{
		ProviderUid: verifyResp.Sub,
		Name:        verifyResp.Name,
	}

	if verifyResp.Email != "" {
		profile.Email = &verifyResp.Email
	}

	return profile, nil
}

func MockValidateLineIdToken(idToken string) (*customer.CustomerProfile, error) {
	if idToken == "" {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthLineTokenInvalid)
	}

	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, errorCodes.NewServiceError(errorCodes.AuthLineTokenInvalid, "invalid token structure", nil)
	}

	if strings.Contains(idToken, "expired") {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthLineTokenExpired)
	}

	if strings.Contains(idToken, "invalid") {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthLineTokenInvalid)
	}

	profile := &customer.CustomerProfile{
		ProviderUid: "U12345678",
		Name:        "Test User",
	}

	if strings.Contains(idToken, "with-email") {
		email := "test@example.com"
		profile.Email = &email
	}

	return profile, nil
}

func (v *LineValidator) SetHTTPClient(client *http.Client) {
	v.httpClient = client
}

func (v *LineValidator) SetVerifyEndpoint(endpoint string) {
	v.verifyEndpoint = endpoint
}