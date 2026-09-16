package arg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// IsUnsupportedLogicalTableError identifies Resource Graph failures that mean a
// recommendation references a logical table unavailable in the current scope/tenant.
// These errors are warning/skip conditions rather than assessment-fatal failures.
func IsUnsupportedLogicalTableError(err error) bool {
	if err == nil {
		return false
	}

	var responseError *azcore.ResponseError
	if !errors.As(err, &responseError) {
		return unsupportedLogicalTableMessage(err.Error())
	}

	if strings.EqualFold(responseError.ErrorCode, "DisallowedLogicalTableName") {
		return true
	}
	if unsupportedLogicalTableMessage(responseError.Error()) {
		return true
	}
	if responseError.RawResponse == nil || responseError.RawResponse.Body == nil {
		return false
	}

	body, readErr := io.ReadAll(responseError.RawResponse.Body)
	if readErr != nil {
		return false
	}
	responseError.RawResponse.Body = io.NopCloser(bytes.NewReader(body))

	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details []struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}
	if strings.EqualFold(payload.Error.Code, "DisallowedLogicalTableName") {
		return true
	}
	for _, detail := range payload.Error.Details {
		if strings.EqualFold(detail.Code, "DisallowedLogicalTableName") || unsupportedLogicalTableMessage(detail.Message) {
			return true
		}
	}
	return false
}

func unsupportedLogicalTableMessage(message string) bool {
	message = strings.ToLower(message)
	if strings.Contains(message, "disallowedlogicaltablename") {
		return true
	}
	return strings.Contains(message, "invalid, unsupported or disallowed") && strings.Contains(message, "logical table")
}
