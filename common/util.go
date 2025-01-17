package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func ParseErrorResponse(endpoint string, url string, body []byte, statusCode int) error {
	var apiError APIError

	if err := json.Unmarshal(body, &apiError); err != nil {
		return fmt.Errorf("%s: error unmarshalling response body: status: %d", endpoint, statusCode)
	}

	errorFormat := "%s url: %s: status: %d, type: %s"
	errorArgs := []interface{}{endpoint, url, statusCode, apiError.Type}

	if apiError.Subtype != "" {
		errorFormat += ", subtype: %s"
		errorArgs = append(errorArgs, apiError.Subtype)
	}

	errorFormat += ", message: %s"
	errorArgs = append(errorArgs, apiError.Message)

	if apiError.Detail != "" {
		errorFormat += ", detail: %s"
		errorArgs = append(errorArgs, apiError.Detail)
	}

	return fmt.Errorf(errorFormat, errorArgs...)
}

func SetUserAgent(userAgent string) string {
	if userAgent == "" {
		return "gocommerce/default-client"
	} else {
		return userAgent
	}
}

func ValidateQueryParams(params QueryParams) error {
	if params.Cursor != "" {
		if params.Filter != "" || params.ModifiedAfter != "" || params.ModifiedBefore != "" ||
			params.SortDirection != "" || params.SortField != "" || params.Status != "" {
			return fmt.Errorf("cannot use cursor alongside other query parameters")
		}
	} else {
		if params.ModifiedAfter != "" && params.ModifiedBefore == "" || params.ModifiedAfter == "" && params.ModifiedBefore != "" {
			return fmt.Errorf("modifiedAfter and modifiedBefore must both be specified together or not at all")
		}
		if params.ModifiedAfter != "" {
			if _, err := time.Parse(time.RFC3339, params.ModifiedAfter); err != nil {
				return fmt.Errorf("modifiedAfter is not a valid ISO 8601 UTC date-time string: %w", err)
			}
		}
		if params.ModifiedBefore != "" {
			if _, err := time.Parse(time.RFC3339, params.ModifiedBefore); err != nil {
				return fmt.Errorf("modifiedBefore is not a valid ISO 8601 UTC date-time string: %w", err)
			}
		}
		if params.Type != "" {
			if err := validateTypeParam(params.Type); err != nil {
				return fmt.Errorf("invalid type: %w", err)
			}
		}
	}

	return nil
}

func validateTypeParam(productType string) error {
	types := strings.Split(productType, ",")
	validTypes := make(map[string]bool)

	for _, t := range types {
		t = strings.TrimSpace(t)
		if t != ProductTypePhysical && t != ProductTypeDigital {
			return fmt.Errorf("type must be either 'PHYSICAL' or 'DIGITAL' (or both comma-separated), got: %s", productType)
		}
		validTypes[t] = true
	}

	if len(validTypes) != len(types) {
		return fmt.Errorf("duplicate types found in: %s", productType)
	}

	return nil
}

// BuildBaseURL constructs the appropriate base URL for API requests.
// During tests, it uses the config.BaseURL if provided, otherwise defaults
// to the Squarespace API URL.
func BuildBaseURL(config *Config, version, path string) (string, error) {
	if config.BaseURL != "" {
		return fmt.Sprintf("%s/%s/%s", config.BaseURL, version, path), nil
	}
	return fmt.Sprintf("https://api.squarespace.com/%s/%s", version, path), nil
}
