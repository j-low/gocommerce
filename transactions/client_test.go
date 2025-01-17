package transactions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/j-low/gocommerce/common"
)

func TestRetrieveAllTransactions(t *testing.T) {
	validTime := time.Now().UTC().Format(time.RFC3339)

	tests := []struct {
		name        string
		params      common.QueryParams
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
		checkQuery  func(*testing.T, string)
	}{
		{
			name:       "successful retrieval",
			params:     common.QueryParams{},
			mockStatus: http.StatusOK,
			mockResp: `{
				"documents": [
					{
						"id": "123",
						"createdOn": "2024-01-01T00:00:00Z",
						"modifiedOn": "2024-01-01T00:00:00Z",
						"total": {"value": "10.00", "currency": "USD"}
					}
				],
				"pagination": {"nextCursor": "next"}
			}`,
			wantErr: false,
		},
		{
			name: "successful retrieval with cursor",
			params: common.QueryParams{
				Cursor: "next-page",
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"documents": [
					{
						"id": "123",
						"createdOn": "2024-01-01T00:00:00Z",
						"modifiedOn": "2024-01-01T00:00:00Z"
					}
				],
				"pagination": {}
			}`,
			wantErr: false,
			checkQuery: func(t *testing.T, query string) {
				if !strings.Contains(query, "cursor=next-page") {
					t.Errorf("expected query to contain cursor parameter, got %s", query)
				}
			},
		},
		{
			name: "successful retrieval with modified dates",
			params: common.QueryParams{
				ModifiedAfter:  validTime,
				ModifiedBefore: validTime,
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"documents": [],
				"pagination": {}
			}`,
			wantErr: false,
			checkQuery: func(t *testing.T, query string) {
				if !strings.Contains(query, "modifiedAfter=") || !strings.Contains(query, "modifiedBefore=") {
					t.Errorf("expected query to contain modified date parameters, got %s", query)
				}
			},
		},
		{
			name: "invalid query parameters",
			params: common.QueryParams{
				Cursor:         "abc",
				ModifiedAfter:  validTime,
				ModifiedBefore: validTime,
			},
			wantErr:     true,
			errContains: "invalid query parameters: cannot use cursor alongside other query parameters",
		},
		{
			name:        "server error",
			params:      common.QueryParams{},
			mockStatus:  http.StatusInternalServerError,
			mockResp:    `{"type":"ERROR","message":"Internal Server Error"}`,
			wantErr:     true,
			errContains: "Internal Server Error",
		},
		{
			name:        "invalid response body",
			params:      common.QueryParams{},
			mockStatus:  http.StatusOK,
			mockResp:    `invalid json`,
			wantErr:     true,
			errContains: "failed to unmarshal response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %s", auth)
				}

				if tt.checkQuery != nil {
					tt.checkQuery(t, r.URL.RawQuery)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			config := &common.Config{
				APIKey:    "test-key",
				Client:    server.Client(),
				UserAgent: "test-agent",
				BaseURL:   server.URL,
			}

			resp, err := RetrieveAllTransactions(context.Background(), config, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAllTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if !tt.wantErr && resp == nil {
				t.Error("expected non-nil response when no error")
			}
		})
	}
}

func TestRetrieveSpecificTransactions(t *testing.T) {
	tests := []struct {
		name           string
		transactionIDs []string
		mockStatus     int
		mockResp       string
		wantErr        bool
		errContains    string
		checkPath      func(*testing.T, string)
	}{
		{
			name:           "successful retrieval",
			transactionIDs: []string{"123", "456"},
			mockStatus:     http.StatusOK,
			mockResp: `{
				"documents": [
					{
						"id": "123",
						"createdOn": "2024-01-01T00:00:00Z",
						"modifiedOn": "2024-01-01T00:00:00Z",
						"total": {"value": "10.00", "currency": "USD"}
					},
					{
						"id": "456",
						"createdOn": "2024-01-01T00:00:00Z",
						"modifiedOn": "2024-01-01T00:00:00Z",
						"total": {"value": "20.00", "currency": "USD"}
					}
				]
			}`,
			wantErr: false,
			checkPath: func(t *testing.T, path string) {
				expected := "123,456"
				if !strings.Contains(path, expected) {
					t.Errorf("expected path to contain %q, got %s", expected, path)
				}
			},
		},
		{
			name:           "empty transaction IDs",
			transactionIDs: []string{},
			wantErr:        true,
			errContains:    "transactionIDs cannot be empty",
		},
		{
			name:           "too many transaction IDs",
			transactionIDs: make([]string, 51),
			wantErr:        true,
			errContains:    "transactionIDs cannot exceed 50 IDs",
		},
		{
			name:           "server error",
			transactionIDs: []string{"123"},
			mockStatus:     http.StatusInternalServerError,
			mockResp:       `{"type":"ERROR","message":"Internal Server Error"}`,
			wantErr:        true,
			errContains:    "Internal Server Error",
		},
		{
			name:           "invalid response body",
			transactionIDs: []string{"123"},
			mockStatus:     http.StatusOK,
			mockResp:       `invalid json`,
			wantErr:        true,
			errContains:    "failed to unmarshal response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %s", auth)
				}

				if tt.checkPath != nil {
					tt.checkPath(t, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			config := &common.Config{
				APIKey:    "test-key",
				Client:    server.Client(),
				UserAgent: "test-agent",
				BaseURL:   server.URL,
			}

			resp, err := RetrieveSpecificTransactions(context.Background(), config, tt.transactionIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpecificTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if !tt.wantErr && resp == nil {
				t.Error("expected non-nil response when no error")
			}
		})
	}
}
