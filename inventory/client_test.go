package inventory

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/j-low/gocommerce/common"
)

func TestRetrieveAllInventory(t *testing.T) {
	tests := []struct {
		name        string
		params      common.QueryParams
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful retrieval",
			params:     common.QueryParams{},
			mockStatus: http.StatusOK,
			mockResp:   `{"inventory": [{"id": "123", "quantity": 5}], "pagination": {"nextCursor": "next"}}`,
			wantErr:    false,
		},
		{
			name: "invalid cursor with other params",
			params: common.QueryParams{
				Cursor: "abc",
				Filter: "test",
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
				// Verify request
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %s", auth)
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

			resp, err := RetrieveAllInventory(context.Background(), config, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAllInventory() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveSpecificInventory(t *testing.T) {
	tests := []struct {
		name         string
		inventoryIDs []string
		mockStatus   int
		mockResp     string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "successful retrieval",
			inventoryIDs: []string{"123", "456"},
			mockStatus:   http.StatusOK,
			mockResp:     `{"inventory": [{"id": "123", "quantity": 5}, {"id": "456", "quantity": 10}]}`,
			wantErr:      false,
		},
		{
			name:         "empty inventory IDs",
			inventoryIDs: []string{},
			wantErr:      true,
			errContains:  "no inventory IDs provided",
		},
		{
			name:         "too many inventory IDs",
			inventoryIDs: make([]string, 51),
			wantErr:      true,
			errContains:  "cannot retrieve more than 50 inventory IDs",
		},
		{
			name:         "server error",
			inventoryIDs: []string{"123"},
			mockStatus:   http.StatusInternalServerError,
			mockResp:     `{"type":"ERROR","message":"Internal Server Error"}`,
			wantErr:      true,
			errContains:  "Internal Server Error",
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

			resp, err := RetrieveSpecificInventory(context.Background(), config, tt.inventoryIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpecificInventory() error = %v, wantErr %v", err, tt.wantErr)
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

func TestAdjustStockQuantities(t *testing.T) {
	tests := []struct {
		name        string
		request     AdjustStockQuantitiesRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful adjustment",
			request: AdjustStockQuantitiesRequest{
				IncrementOperations: []QuantityOperation{
					{VariantID: "123", Quantity: 5},
				},
			},
			mockStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name: "server error",
			request: AdjustStockQuantitiesRequest{
				IncrementOperations: []QuantityOperation{
					{VariantID: "123", Quantity: 5},
				},
			},
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"ERROR","message":"Invalid request"}`,
			wantErr:     true,
			errContains: "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %s", auth)
				}
				if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
					t.Errorf("expected Content-Type header 'application/json', got %s", contentType)
				}

				w.WriteHeader(tt.mockStatus)
				if tt.mockResp != "" {
					w.Write([]byte(tt.mockResp))
				}
			}))
			defer server.Close()

			config := &common.Config{
				APIKey:    "test-key",
				Client:    server.Client(),
				UserAgent: "test-agent",
				BaseURL:   server.URL,
			}

			status, err := AdjustStockQuantities(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdjustStockQuantities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if !tt.wantErr && status != tt.mockStatus {
				t.Errorf("expected status %d, got %d", tt.mockStatus, status)
			}
		})
	}
}
