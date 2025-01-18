package webhooks

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/j-low/gocommerce/common"
)

func TestCreateWebhookSubscription(t *testing.T) {
	tests := []struct {
		name        string
		config      *common.Config
		request     WebhookSubscriptionRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation",
			config: &common.Config{
				AccessToken: "test-token",
			},
			request: WebhookSubscriptionRequest{
				EndpointURL: "https://example.com/webhook",
				Topics:      []string{"order.created"},
			},
			mockStatus: http.StatusCreated,
			mockResp: `{
				"id": "webhook-123",
				"endpointUrl": "https://example.com/webhook",
				"topics": ["order.created"],
				"secret": "secret123",
				"createdOn": "2024-01-01T00:00:00Z",
				"updatedOn": "2024-01-01T00:00:00Z"
			}`,
		},
		{
			name: "missing access token",
			config: &common.Config{
				AccessToken: "",
			},
			request: WebhookSubscriptionRequest{
				EndpointURL: "https://example.com/webhook",
				Topics:      []string{"order.created"},
			},
			wantErr:     true,
			errContains: "access token is required",
		},
		{
			name: "empty topics",
			config: &common.Config{
				AccessToken: "test-token",
			},
			request: WebhookSubscriptionRequest{
				EndpointURL: "https://example.com/webhook",
				Topics:      []string{},
			},
			wantErr:     true,
			errContains: "topics cannot be empty",
		},
		{
			name: "server error",
			config: &common.Config{
				AccessToken: "test-token",
			},
			request: WebhookSubscriptionRequest{
				EndpointURL: "https://example.com/webhook",
				Topics:      []string{"order.created"},
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
				if auth := r.Header.Get("Authorization"); auth != "Bearer "+tt.config.AccessToken {
					t.Errorf("expected Authorization header 'Bearer %s', got %s", tt.config.AccessToken, auth)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			tt.config.Client = server.Client()
			tt.config.BaseURL = server.URL

			resp, err := CreateWebhookSubscription(context.Background(), tt.config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWebhookSubscription() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUpdateWebhookSubscription(t *testing.T) {
	tests := []struct {
		name           string
		config         *common.Config
		subscriptionID string
		request        WebhookSubscriptionRequest
		mockStatus     int
		mockResp       string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful update",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			request: WebhookSubscriptionRequest{
				EndpointURL: "https://example.com/webhook-updated",
				Topics:      []string{"order.created", "order.fulfilled"},
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"id": "webhook-123",
				"endpointUrl": "https://example.com/webhook-updated",
				"topics": ["order.created", "order.fulfilled"]
			}`,
		},
		{
			name: "missing access token",
			config: &common.Config{
				AccessToken: "",
			},
			subscriptionID: "webhook-123",
			wantErr:        true,
			errContains:    "access token is required",
		},
		{
			name: "empty subscription ID",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name: "empty topics array",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			request: WebhookSubscriptionRequest{
				EndpointURL: "https://example.com/webhook",
				Topics:      []string{},
			},
			wantErr:     true,
			errContains: "topics cannot be an empty array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/1.0/webhook_subscriptions/" + tt.subscriptionID
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			tt.config.Client = server.Client()
			tt.config.BaseURL = server.URL

			resp, err := UpdateWebhookSubscription(context.Background(), tt.config, tt.subscriptionID, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateWebhookSubscription() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveAllWebhookSubscriptions(t *testing.T) {
	tests := []struct {
		name        string
		config      *common.Config
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful retrieval",
			config: &common.Config{
				AccessToken: "test-token",
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"webhookSubscriptions": [
					{
						"id": "webhook-123",
						"endpointUrl": "https://example.com/webhook",
						"topics": ["order.created"]
					}
				]
			}`,
		},
		{
			name: "missing access token",
			config: &common.Config{
				AccessToken: "",
			},
			wantErr:     true,
			errContains: "access token is required",
		},
		{
			name: "server error",
			config: &common.Config{
				AccessToken: "test-token",
			},
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

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			tt.config.Client = server.Client()
			tt.config.BaseURL = server.URL

			resp, err := RetrieveAllWebhookSubscriptions(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAllWebhookSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveSpecificWebhookSubscription(t *testing.T) {
	tests := []struct {
		name           string
		config         *common.Config
		subscriptionID string
		mockStatus     int
		mockResp       string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful retrieval",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			mockStatus:     http.StatusOK,
			mockResp: `{
				"id": "webhook-123",
				"endpointUrl": "https://example.com/webhook",
				"topics": ["order.created"]
			}`,
		},
		{
			name: "missing access token",
			config: &common.Config{
				AccessToken: "",
			},
			subscriptionID: "webhook-123",
			wantErr:        true,
			errContains:    "access token is required",
		},
		{
			name: "empty subscription ID",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name: "not found",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			mockStatus:     http.StatusNotFound,
			mockResp:       `{"type":"ERROR","message":"Webhook subscription not found"}`,
			wantErr:        true,
			errContains:    "Webhook subscription not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				expectedPath := "/1.0/webhook_subscriptions/" + tt.subscriptionID
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			tt.config.Client = server.Client()
			tt.config.BaseURL = server.URL

			resp, err := RetrieveSpecificWebhookSubscription(context.Background(), tt.config, tt.subscriptionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpecificWebhookSubscription() error = %v, wantErr %v", err, tt.wantErr)
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

func TestDeleteWebhookSubscription(t *testing.T) {
	tests := []struct {
		name           string
		config         *common.Config
		subscriptionID string
		mockStatus     int
		mockResp       string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful deletion",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			mockStatus:     http.StatusNoContent,
		},
		{
			name: "missing access token",
			config: &common.Config{
				AccessToken: "",
			},
			subscriptionID: "webhook-123",
			wantErr:        true,
			errContains:    "access token is required",
		},
		{
			name: "empty subscription ID",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name: "not found",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			mockStatus:     http.StatusNotFound,
			mockResp:       `{"type":"ERROR","message":"Webhook subscription not found"}`,
			wantErr:        true,
			errContains:    "Webhook subscription not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}

				w.WriteHeader(tt.mockStatus)
				if tt.mockResp != "" {
					w.Write([]byte(tt.mockResp))
				}
			}))
			defer server.Close()

			tt.config.Client = server.Client()
			tt.config.BaseURL = server.URL

			status, err := DeleteWebhookSubscription(context.Background(), tt.config, tt.subscriptionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWebhookSubscription() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSendTestNotification(t *testing.T) {
	tests := []struct {
		name           string
		config         *common.Config
		subscriptionID string
		request        SendTestNotificationRequest
		mockStatus     int
		mockResp       string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful test notification",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			request: SendTestNotificationRequest{
				Topic: "order.created",
			},
			mockStatus: http.StatusOK,
			mockResp:   `{"statusCode": 200}`,
		},
		{
			name: "missing access token",
			config: &common.Config{
				AccessToken: "",
			},
			subscriptionID: "webhook-123",
			request: SendTestNotificationRequest{
				Topic: "order.created",
			},
			wantErr:     true,
			errContains: "access token is required",
		},
		{
			name: "empty subscription ID",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "",
			request: SendTestNotificationRequest{
				Topic: "order.created",
			},
			wantErr:     true,
			errContains: "subscriptionID cannot be empty",
		},
		{
			name: "empty topic",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			request: SendTestNotificationRequest{
				Topic: "",
			},
			wantErr:     true,
			errContains: "topic is required",
		},
		{
			name: "server error",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			request: SendTestNotificationRequest{
				Topic: "order.created",
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
				expectedPath := fmt.Sprintf("/1.0/webhook_subscriptions/%s/actions/sendTestNotification", tt.subscriptionID)
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			tt.config.Client = server.Client()
			tt.config.BaseURL = server.URL

			resp, err := SendTestNotification(context.Background(), tt.config, tt.subscriptionID, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendTestNotification() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRotateSubscriptionSecret(t *testing.T) {
	tests := []struct {
		name           string
		config         *common.Config
		subscriptionID string
		mockStatus     int
		mockResp       string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful secret rotation",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			mockStatus:     http.StatusOK,
			mockResp:       `{"secret": "new-secret-123"}`,
		},
		{
			name: "missing access token",
			config: &common.Config{
				AccessToken: "",
			},
			subscriptionID: "webhook-123",
			wantErr:        true,
			errContains:    "access token is required",
		},
		{
			name: "empty subscription ID",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name: "server error",
			config: &common.Config{
				AccessToken: "test-token",
			},
			subscriptionID: "webhook-123",
			mockStatus:     http.StatusBadRequest,
			mockResp:       `{"type":"ERROR","message":"Invalid request"}`,
			wantErr:        true,
			errContains:    "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/1.0/webhook_subscriptions/%s/actions/rotateSecret", tt.subscriptionID)
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			tt.config.Client = server.Client()
			tt.config.BaseURL = server.URL

			resp, err := RotateSubscriptionSecret(context.Background(), tt.config, tt.subscriptionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RotateSubscriptionSecret() error = %v, wantErr %v", err, tt.wantErr)
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
