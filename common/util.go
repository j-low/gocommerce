package common

import (
	"encoding/json"
	"fmt"
)

func ParseErrorResponse(body []byte, statusCode int) error {
	var apiError APIError

	if err := json.Unmarshal(body, &apiError); err != nil {
		return fmt.Errorf("unexpected status code: %d, body: %s", statusCode, string(body))
	}

	errorMessage := fmt.Sprintf("status: %d, type: %s, subtype: %s, message: %s, detail: %s",
		statusCode, apiError.Type, apiError.Subtype, apiError.Message, apiError.Detail)

	return fmt.Errorf(errorMessage)
}

func SetUserAgent(userAgent string) string {
	if userAgent == "" {
		return "gocommerce/default-client"
	} else {
		return userAgent
	}
}
