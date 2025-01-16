package common

import (
	"encoding/json"
	"fmt"
	"time"
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

func ValidateQueryParams(queryParams QueryParams) error {
	if queryParams.Cursor != "" {
		if queryParams.Filter != "" || queryParams.ModifiedAfter != "" || queryParams.ModifiedBefore != "" ||
			queryParams.SortDirection != "" || queryParams.SortField != "" || queryParams.Status != "" {
			return fmt.Errorf("cannot use cursor alongside other query parameters")
		}
	} else {
		if err := validateModifiedBeforeAfterQueryParams(queryParams.ModifiedAfter, queryParams.ModifiedBefore); err != nil {
			return fmt.Errorf("invalid modifiedAfter or modifiedBefore: %w", err)
		}
	}

	return nil
}

func validateModifiedBeforeAfterQueryParams(modifiedAfter, modifiedBefore string) error {
	if modifiedAfter != "" && modifiedBefore == "" || modifiedAfter == "" && modifiedBefore != "" {
		return fmt.Errorf("modifiedAfter and modifiedBefore must both be specified together or not at all")
	}
	if modifiedAfter != "" {
		if _, err := time.Parse(time.RFC3339, modifiedAfter); err != nil {
			return fmt.Errorf("modifiedAfter is not a valid ISO 8601 UTC date-time string: %w", err)
		}
	}
	if modifiedBefore != "" {
		if _, err := time.Parse(time.RFC3339, modifiedBefore); err != nil {
			return fmt.Errorf("modifiedBefore is not a valid ISO 8601 UTC date-time string: %w", err)
		}
	}

	return nil
}
