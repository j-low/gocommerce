package common

import (
	"strings"
	"testing"
	"time"
)

func TestParseErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		url        string
		statusCode int
		body       []byte
		wantErr    string
	}{
		{
			name:       "basic error",
			endpoint:   "TestEndpoint",
			url:        "http://example.com/api",
			statusCode: 400,
			body:       []byte(`{"type":"ERROR_TYPE","message":"Error occurred"}`),
			wantErr:    "TestEndpoint url: http://example.com/api: status: 400, type: ERROR_TYPE, message: Error occurred",
		},
		{
			name:       "error with subtype",
			endpoint:   "TestEndpoint",
			url:        "http://example.com/api",
			statusCode: 400,
			body:       []byte(`{"type":"ERROR_TYPE","subtype":"SUB_ERROR","message":"Error occurred"}`),
			wantErr:    "TestEndpoint url: http://example.com/api: status: 400, type: ERROR_TYPE, subtype: SUB_ERROR, message: Error occurred",
		},
		{
			name:       "error with detail",
			endpoint:   "TestEndpoint",
			url:        "http://example.com/api",
			statusCode: 400,
			body:       []byte(`{"type":"ERROR_TYPE","message":"Error occurred","detail":"Additional details"}`),
			wantErr:    "TestEndpoint url: http://example.com/api: status: 400, type: ERROR_TYPE, message: Error occurred, detail: Additional details",
		},
		{
			name:       "invalid json",
			endpoint:   "TestEndpoint",
			url:        "http://example.com/api",
			statusCode: 400,
			body:       []byte(`invalid json`),
			wantErr:    "TestEndpoint: error unmarshalling response body: status: 400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseErrorResponse(tt.endpoint, tt.url, tt.body, tt.statusCode)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("ParseErrorResponse() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestSetUserAgent(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		want      string
	}{
		{
			name:      "empty user agent",
			userAgent: "",
			want:      "gocommerce/default-client",
		},
		{
			name:      "custom user agent",
			userAgent: "custom-agent",
			want:      "custom-agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetUserAgent(tt.userAgent); got != tt.want {
				t.Errorf("SetUserAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateQueryParams(t *testing.T) {
	validTime := time.Now().UTC().Format(time.RFC3339)

	tests := []struct {
		name    string
		params  QueryParams
		wantErr string
	}{
		{
			name:   "valid empty params",
			params: QueryParams{},
		},
		{
			name: "cursor with other params",
			params: QueryParams{
				Cursor: "abc123",
				Filter: "some-filter",
			},
			wantErr: "cannot use cursor alongside other query parameters",
		},
		{
			name: "only modifiedAfter",
			params: QueryParams{
				ModifiedAfter: validTime,
			},
			wantErr: "modifiedAfter and modifiedBefore must both be specified together or not at all",
		},
		{
			name: "only modifiedBefore",
			params: QueryParams{
				ModifiedBefore: validTime,
			},
			wantErr: "modifiedAfter and modifiedBefore must both be specified together or not at all",
		},
		{
			name: "invalid modifiedAfter format",
			params: QueryParams{
				ModifiedAfter:  "invalid-date",
				ModifiedBefore: validTime,
			},
			wantErr: "modifiedAfter is not a valid ISO 8601 UTC date-time string",
		},
		{
			name: "invalid modifiedBefore format",
			params: QueryParams{
				ModifiedAfter:  validTime,
				ModifiedBefore: "invalid-date",
			},
			wantErr: "modifiedBefore is not a valid ISO 8601 UTC date-time string",
		},
		{
			name: "valid modified dates",
			params: QueryParams{
				ModifiedAfter:  validTime,
				ModifiedBefore: validTime,
			},
		},
		{
			name: "invalid type",
			params: QueryParams{
				Type: "INVALID",
			},
			wantErr: "invalid type: type must be either 'PHYSICAL' or 'DIGITAL' (or both comma-separated), got: INVALID",
		},
		{
			name: "valid single type",
			params: QueryParams{
				Type: "PHYSICAL",
			},
		},
		{
			name: "valid multiple types",
			params: QueryParams{
				Type: "PHYSICAL,DIGITAL",
			},
		},
		{
			name: "duplicate types",
			params: QueryParams{
				Type: "PHYSICAL,PHYSICAL",
			},
			wantErr: "invalid type: duplicate types found in: PHYSICAL,PHYSICAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQueryParams(tt.params)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("ValidateQueryParams() unexpected error = %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("ValidateQueryParams() error = %v, want to contain %v", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestValidateTypeParam(t *testing.T) {
	tests := []struct {
		name      string
		typeParam string
		wantErr   string
	}{
		{
			name:      "valid single type - PHYSICAL",
			typeParam: "PHYSICAL",
		},
		{
			name:      "valid single type - DIGITAL",
			typeParam: "DIGITAL",
		},
		{
			name:      "valid multiple types",
			typeParam: "PHYSICAL,DIGITAL",
		},
		{
			name:      "valid multiple types with spaces",
			typeParam: "PHYSICAL, DIGITAL",
		},
		{
			name:      "invalid type",
			typeParam: "INVALID",
			wantErr:   "type must be either 'PHYSICAL' or 'DIGITAL'",
		},
		{
			name:      "duplicate types",
			typeParam: "PHYSICAL,PHYSICAL",
			wantErr:   "duplicate types found in: PHYSICAL,PHYSICAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTypeParam(tt.typeParam)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("validateTypeParam() unexpected error = %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("validateTypeParam() error = %v, want to contain %v", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestBuildBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		version string
		path    string
		want    string
		wantErr bool
	}{
		{
			name: "default squarespace URL",
			config: &Config{
				APIKey:    "test-key",
				UserAgent: "test-agent",
			},
			version: "1.0",
			path:    "commerce/inventory",
			want:    "https://api.squarespace.com/1.0/commerce/inventory",
			wantErr: false,
		},
		{
			name: "custom base URL",
			config: &Config{
				APIKey:    "test-key",
				UserAgent: "test-agent",
				BaseURL:   "http://localhost:8080",
			},
			version: "1.0",
			path:    "commerce/inventory",
			want:    "http://localhost:8080/1.0/commerce/inventory",
			wantErr: false,
		},
		{
			name: "empty version",
			config: &Config{
				APIKey:    "test-key",
				UserAgent: "test-agent",
			},
			version: "",
			path:    "commerce/inventory",
			want:    "https://api.squarespace.com//commerce/inventory",
			wantErr: false,
		},
		{
			name: "empty path",
			config: &Config{
				APIKey:    "test-key",
				UserAgent: "test-agent",
			},
			version: "1.0",
			path:    "",
			want:    "https://api.squarespace.com/1.0/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildBaseURL(tt.config, tt.version, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildBaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
