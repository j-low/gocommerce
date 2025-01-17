package orders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/j-low/gocommerce/common"
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name        string
		request     CreateOrderRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation",
			request: CreateOrderRequest{
				ChannelName:            "API",
				ExternalOrderReference: "test-123",
				CustomerEmail:          "test@example.com",
				LineItems: []LineItem{
					{
						LineItemType:  "PHYSICAL",
						VariantID:     "123",
						Quantity:      1,
						UnitPricePaid: common.Amount{Value: "10.00", Currency: "USD"},
					},
				},
				PriceTaxInterpretation: "INCLUSIVE",
				GrandTotal:             common.Amount{Value: "10.00", Currency: "USD"},
				CreatedOn:              "2024-01-01T00:00:00Z",
			},
			mockStatus: http.StatusCreated,
			mockResp:   `{"id": "order-123", "orderNumber": "123"}`,
			wantErr:    false,
		},
		{
			name: "invalid request",
			request: CreateOrderRequest{
				// Missing required fields
				ChannelName: "API",
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
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			config := &common.Config{
				APIKey:    "test-key",
				Client:    server.Client(),
				UserAgent: "test-agent",
				BaseURL:   server.URL,
			}

			resp, err := CreateOrder(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
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

func TestFulfillOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		request     FulfillOrderRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful fulfillment",
			orderID: "order-123",
			request: FulfillOrderRequest{
				ShouldSendNotification: true,
				Shipments: []Shipment{
					{
						ShipDate:       "2024-01-01",
						CarrierName:    "UPS",
						Service:        "Ground",
						TrackingNumber: "1Z999999999",
					},
				},
			},
			mockStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:    "invalid order ID",
			orderID: "invalid-id",
			request: FulfillOrderRequest{
				ShouldSendNotification: true,
				Shipments:              []Shipment{},
			},
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"ERROR","message":"Invalid order ID"}`,
			wantErr:     true,
			errContains: "Invalid order ID",
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

			status, err := FulfillOrder(context.Background(), config, tt.orderID, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("FulfillOrder() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveAllOrders(t *testing.T) {
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
			mockResp: `{
				"result": [
					{
						"id": "order-123",
						"orderNumber": "123",
						"customerEmail": "test@example.com"
					}
				],
				"pagination": {"nextCursor": "next"}
			}`,
			wantErr: false,
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

			resp, err := RetrieveAllOrders(context.Background(), config, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAllOrders() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveSpecificOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful retrieval",
			orderID:    "order-123",
			mockStatus: http.StatusOK,
			mockResp: `{
				"id": "order-123",
				"orderNumber": "123",
				"customerEmail": "test@example.com",
				"lineItems": [
					{
						"id": "item-1",
						"variantId": "variant-1",
						"quantity": 1,
						"unitPricePaid": {"value": "10.00", "currency": "USD"}
					}
				],
				"grandTotal": {"value": "10.00", "currency": "USD"}
			}`,
			wantErr: false,
		},
		{
			name:        "order not found",
			orderID:     "non-existent",
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Order not found"}`,
			wantErr:     true,
			errContains: "Order not found",
		},
		{
			name:        "empty order ID",
			orderID:     "",
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"ERROR","message":"Order ID is required"}`,
			wantErr:     true,
			errContains: "Order ID is required",
		},
		{
			name:        "server error",
			orderID:     "order-123",
			mockStatus:  http.StatusInternalServerError,
			mockResp:    `{"type":"ERROR","message":"Internal Server Error"}`,
			wantErr:     true,
			errContains: "Internal Server Error",
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

				expectedPath := "/1.0/commerce/orders/" + tt.orderID
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
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

			resp, err := RetrieveSpecificOrder(context.Background(), config, tt.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpecificOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Error("expected non-nil response when no error")
					return
				}

				if resp.ID != "order-123" && tt.name == "successful retrieval" {
					t.Errorf("expected order ID 'order-123', got %s", resp.ID)
				}
			}
		})
	}
}
