package products

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/j-low/gocommerce/common"
)

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name        string
		request     CreateProductRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation",
			request: CreateProductRequest{
				Type:        "PHYSICAL",
				StorePageID: "store-page-123",
				Name:        "Test Product",
				Variants: []ProductVariant{
					{
						SKU: "TEST-123",
						Pricing: Pricing{
							BasePrice: common.Amount{Value: "10.00", Currency: "USD"},
						},
						Stock: Stock{
							Quantity:  100,
							Unlimited: false,
						},
					},
				},
			},
			mockStatus: http.StatusCreated,
			mockResp: `{
                "id": "product-123",
                "type": "PHYSICAL",
                "name": "Test Product"
            }`,
			wantErr: false,
		},
		{
			name: "invalid request",
			request: CreateProductRequest{
				// Missing required fields
				Name: "Test Product",
			},
			mockStatus: http.StatusBadRequest,
			mockResp: `{
                "type": "ERROR",
                "message": "Invalid request: missing required fields"
            }`,
			wantErr:     true,
			errContains: "Invalid request",
		},
		{
			name: "server error",
			request: CreateProductRequest{
				Type: "PHYSICAL",
				Name: "Test Product",
			},
			mockStatus: http.StatusInternalServerError,
			mockResp: `{
                "type": "ERROR",
                "message": "Internal server error"
            }`,
			wantErr:     true,
			errContains: "Internal server error",
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

			resp, err := CreateProduct(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
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

func TestCreateProductVariant(t *testing.T) {
	tests := []struct {
		name        string
		request     CreateProductVariantRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation",
			request: CreateProductVariantRequest{
				ProductID: "product-123",
				SKU:       "TEST-VAR-123",
				Pricing: Pricing{
					BasePrice: common.Amount{Value: "29.99", Currency: "USD"},
				},
				Stock: Stock{
					Quantity:  100,
					Unlimited: false,
				},
				ShippingMeasurements: ShippingMeasurements{
					Weight: Weight{
						Unit:  "LB",
						Value: 2.5,
					},
					Dimensions: Dimensions{
						Unit:   "IN",
						Length: 10,
						Width:  5,
						Height: 2,
					},
				},
			},
			mockStatus: http.StatusCreated,
			mockResp: `{
				"id": "variant-123",
				"sku": "TEST-VAR-123",
				"pricing": {
					"basePrice": {"value": "29.99", "currency": "USD"}
				}
			}`,
			wantErr: false,
		},
		{
			name: "missing required fields",
			request: CreateProductVariantRequest{
				ProductID: "product-123",
				// Missing SKU and Pricing
			},
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"ERROR","message":"SKU and Pricing are required"}`,
			wantErr:     true,
			errContains: "SKU and Pricing are required",
		},
		{
			name: "missing product ID",
			request: CreateProductVariantRequest{
				// Missing ProductID
				SKU: "TEST-VAR-123",
				Pricing: Pricing{
					BasePrice: common.Amount{Value: "29.99", Currency: "USD"},
				},
			},
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"ERROR","message":"productID is required"}`,
			wantErr:     true,
			errContains: "productID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/variants", tt.request.ProductID)
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

			resp, err := CreateProductVariant(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProductVariant() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUploadProductImage(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		filename    string
		imageData   []byte
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful upload",
			productID:  "product-123",
			filename:   "test-image.jpg",
			imageData:  []byte("fake-image-data"),
			mockStatus: http.StatusAccepted, // 202 ACCEPTED
			mockResp:   `{"imageId": "5ed539bc8367410cdc0c984a"}`,
			wantErr:    false,
		},
		{
			name:        "product not found",
			productID:   "non-existent",
			filename:    "test.jpg",
			imageData:   []byte("test-data"),
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"INVALID_REQUEST_ERROR","message":"Product not found","subtype":"INVALID_ARGUMENT"}`,
			wantErr:     true,
			errContains: "Product not found",
		},
		{
			name:        "image limit reached",
			productID:   "product-123",
			filename:    "test.jpg",
			imageData:   []byte("test-data"),
			mockStatus:  http.StatusConflict,
			mockResp:    `{"type":"CONFLICT","message":"Product has reached image limit","subtype":"IMAGE_LIMIT_REACHED"}`,
			wantErr:     true,
			errContains: "Product has reached image limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpfile, err := os.CreateTemp("", "test-*.jpg")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write(tt.imageData); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if contentType := r.Header.Get("Content-Type"); !strings.Contains(contentType, "multipart/form-data") {
					t.Errorf("expected Content-Type to contain multipart/form-data, got %s", contentType)
				}

				err := r.ParseMultipartForm(20 << 20) // 20MB max
				if err != nil {
					t.Errorf("failed to parse multipart form: %v", err)
				}

				file, _, err := r.FormFile("file") // Note: form field must be "file"
				if err != nil {
					t.Errorf("failed to get form file: %v", err)
				}
				defer file.Close()

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

			resp, err := UploadProductImage(context.Background(), config, tt.productID, tmpfile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadProductImage() error = %v, wantErr %v", err, tt.wantErr)
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
				if resp.ImageID == "" {
					t.Error("expected non-empty imageId in response")
				}
			}
		})
	}
}

func TestRetrieveAllStorePages(t *testing.T) {
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
				"storePages": [
					{
						"id": "page-123",
						"title": "Test Page",
						"url": "/test-page"
					}
				],
				"pagination": {
					"nextCursor": "next-cursor"
				}
			}`,
			wantErr: false,
		},
		{
			name: "with cursor",
			params: common.QueryParams{
				Cursor: "test-cursor",
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"storePages": [
					{
						"id": "page-456",
						"title": "Another Page",
						"url": "/another-page"
					}
				]
			}`,
			wantErr: false,
		},
		{
			name:       "server error",
			params:     common.QueryParams{},
			mockStatus: http.StatusInternalServerError,
			mockResp: `{
				"type": "ERROR",
				"message": "Internal server error"
			}`,
			wantErr:     true,
			errContains: "Internal server error",
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

				// Verify cursor parameter if present
				if tt.params.Cursor != "" {
					cursor := r.URL.Query().Get("cursor")
					if cursor != tt.params.Cursor {
						t.Errorf("expected cursor %s, got %s", tt.params.Cursor, cursor)
					}
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

			resp, err := RetrieveAllStorePages(context.Background(), config, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAllStorePages() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveAllProducts(t *testing.T) {
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
                "products": [
                    {
                        "id": "product-123",
                        "type": "PHYSICAL",
                        "name": "Test Product"
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
			errContains: "cannot use cursor alongside other query parameters",
		},
		{
			name:       "server error",
			params:     common.QueryParams{},
			mockStatus: http.StatusInternalServerError,
			mockResp: `{
                "type": "ERROR",
                "message": "Internal server error"
            }`,
			wantErr:     true,
			errContains: "Internal server error",
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

			resp, err := RetrieveAllProducts(context.Background(), config, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAllProducts() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveSpecificProducts(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful retrieval",
			productID:  "product-123",
			mockStatus: http.StatusOK,
			mockResp: `{
				"id": "product-123",
				"type": "PHYSICAL",
				"name": "Test Product",
				"variants": [
					{
						"id": "variant-1",
						"sku": "TEST-123",
						"pricing": {
							"basePrice": {"value": "10.00", "currency": "USD"}
						}
					}
				]
			}`,
			wantErr: false,
		},
		{
			name:        "product not found",
			productID:   "non-existent",
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Product not found"}`,
			wantErr:     true,
			errContains: "Product not found",
		},
		{
			name:        "server error",
			productID:   "product-123",
			mockStatus:  http.StatusInternalServerError,
			mockResp:    `{"type":"ERROR","message":"Internal server error"}`,
			wantErr:     true,
			errContains: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				expectedPath := "/1.0/commerce/products/" + tt.productID
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

			resp, err := RetrieveSpecificProducts(context.Background(), config, []string{tt.productID})
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpecificProduct() error = %v, wantErr %v", err, tt.wantErr)
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

func TestGetProductImageUploadStatus(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		imageID     string
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful status retrieval",
			productID:  "product-123",
			imageID:    "image-123",
			mockStatus: http.StatusOK,
			mockResp: `{
				"status": "READY",
				"image": {
					"id": "image-123",
					"url": "https://example.com/image.jpg"
				}
			}`,
			wantErr: false,
		},
		{
			name:        "image not found",
			productID:   "product-123",
			imageID:     "non-existent",
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Image not found"}`,
			wantErr:     true,
			errContains: "Image not found",
		},
		{
			name:        "server error",
			productID:   "product-123",
			imageID:     "image-123",
			mockStatus:  http.StatusInternalServerError,
			mockResp:    `{"type":"ERROR","message":"Internal server error"}`,
			wantErr:     true,
			errContains: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/images/%s/status", tt.productID, tt.imageID)
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

			resp, err := GetProductImageUploadStatus(context.Background(), config, tt.productID, tt.imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductImageUploadStatus() error = %v, wantErr %v", err, tt.wantErr)
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

func TestAssignProductImageToVariant(t *testing.T) {
	tests := []struct {
		name        string
		request     AssignProductImageToVariantRequest
		mockStatus  int
		mockResp    string
		wantStatus  int
		wantErr     bool
		errContains string
	}{
		{
			name: "successful assignment",
			request: AssignProductImageToVariantRequest{
				ProductID: "product-123",
				VariantID: "variant-123",
				ImageID:   "image-123",
			},
			mockStatus: http.StatusNoContent,
			mockResp:   "",
			wantStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name: "variant not found",
			request: AssignProductImageToVariantRequest{
				ProductID: "product-123",
				VariantID: "non-existent",
				ImageID:   "image-123",
			},
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Variant not found"}`,
			wantStatus:  http.StatusNotFound,
			wantErr:     true,
			errContains: "Variant not found",
		},
		{
			name: "invalid request",
			request: AssignProductImageToVariantRequest{
				ProductID: "",
				VariantID: "",
				ImageID:   "",
			},
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"ERROR","message":"Invalid request parameters"}`,
			wantStatus:  http.StatusBadRequest,
			wantErr:     true,
			errContains: "Invalid request parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/variants/%s/image",
					tt.request.ProductID, tt.request.VariantID)
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

			status, err := AssignProductImageToVariant(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignProductImageToVariant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if status != tt.wantStatus {
				t.Errorf("AssignProductImageToVariant() status = %v, want %v", status, tt.wantStatus)
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestReorderProductImage(t *testing.T) {
	imageID_1 := "image-456"
	imageID_2 := "image-789"
	invalidFormat := "invalid-format"

	tests := []struct {
		name        string
		request     ReorderProductImageRequest
		mockStatus  int
		mockResp    string
		wantStatus  int
		wantErr     bool
		errContains string
		checkBody   bool
		wantBody    string
	}{
		{
			name: "successful reorder - after specific image",
			request: ReorderProductImageRequest{
				ProductID:    "product-123",
				ImageID:      imageID_1,
				AfterImageID: &imageID_2,
			},
			mockStatus: http.StatusNoContent,
			wantStatus: http.StatusNoContent,
			wantErr:    false,
			checkBody:  true,
			wantBody:   `{"afterImageId":"image-789"}`,
		},
		{
			name: "successful reorder - move to top",
			request: ReorderProductImageRequest{
				ProductID:    "product-123",
				ImageID:      imageID_1,
				AfterImageID: nil,
			},
			mockStatus: http.StatusNoContent,
			wantStatus: http.StatusNoContent,
			wantErr:    false,
			checkBody:  true,
			wantBody:   `{"afterImageId":null}`,
		},
		{
			name: "image not found",
			request: ReorderProductImageRequest{
				ProductID:    "product-123",
				ImageID:      "non-existent",
				AfterImageID: &imageID_2,
			},
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"INVALID_REQUEST_ERROR","subtype":"INVALID_ARGUMENT","message":"Image not found"}`,
			wantStatus:  http.StatusNotFound,
			wantErr:     true,
			errContains: "Image not found",
		},
		{
			name: "invalid product type",
			request: ReorderProductImageRequest{
				ProductID:    "digital-product-123",
				ImageID:      imageID_1,
				AfterImageID: &imageID_2,
			},
			mockStatus:  http.StatusMethodNotAllowed,
			mockResp:    `{"type":"METHOD_NOT_ALLOWED","subtype":"OPERATION_NOT_ALLOWED_FOR_PRODUCT_TYPE","message":"Operation not allowed for digital products"}`,
			wantStatus:  http.StatusMethodNotAllowed,
			wantErr:     true,
			errContains: "Operation not allowed for digital products",
		},
		{
			name: "invalid request body",
			request: ReorderProductImageRequest{
				ProductID:    "product-123",
				ImageID:      imageID_1,
				AfterImageID: &invalidFormat,
			},
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"INVALID_REQUEST_ERROR","message":"The request body does not conform to the required specification"}`,
			wantStatus:  http.StatusBadRequest,
			wantErr:     true,
			errContains: "The request body does not conform to the required specification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/images/%s/order",
					tt.request.ProductID, tt.request.ImageID)
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
				}

				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %s", auth)
				}
				if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
					t.Errorf("expected Content-Type header 'application/json', got %s", contentType)
				}

				if tt.checkBody {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					if string(body) != tt.wantBody {
						t.Errorf("expected request body %s, got %s", tt.wantBody, string(body))
					}
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

			status, err := ReorderProductImage(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReorderProductImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if status != tt.wantStatus {
				t.Errorf("ReorderProductImage() status = %v, want %v", status, tt.wantStatus)
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		request     UpdateProductRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:      "successful update",
			productID: "product-123",
			request: UpdateProductRequest{
				Name: "Updated Product",
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"id": "product-123",
				"name": "Updated Product"
			}`,
			wantErr: false,
		},
		{
			name:       "invalid request",
			productID:  "product-123",
			request:    UpdateProductRequest{},
			mockStatus: http.StatusBadRequest,
			mockResp: `{
				"type": "ERROR",
				"message": "Invalid request"
			}`,
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
				expectedPath := "/1.0/commerce/products/" + tt.productID
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

			resp, err := UpdateProduct(context.Background(), config, tt.productID, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUpdateProductVariant(t *testing.T) {
	tests := []struct {
		name        string
		request     UpdateProductVariantRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful update",
			request: UpdateProductVariantRequest{
				ProductID: "product-123",
				VariantID: "variant-123",
				SKU:       "TEST-VAR-123-UPDATED",
				Pricing: Pricing{
					BasePrice: common.Amount{Value: "39.99", Currency: "USD"},
					OnSale:    true,
					SalePrice: common.Amount{Value: "29.99", Currency: "USD"},
				},
				ShippingMeasurements: ShippingMeasurements{
					Weight: Weight{
						Unit:  "LB",
						Value: 3.0,
					},
				},
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"id": "variant-123",
				"sku": "TEST-VAR-123-UPDATED",
				"pricing": {
					"basePrice": {"value": "39.99", "currency": "USD"},
					"onSale": true,
					"salePrice": {"value": "29.99", "currency": "USD"}
				}
			}`,
			wantErr: false,
		},
		{
			name: "variant not found",
			request: UpdateProductVariantRequest{
				ProductID: "product-123",
				VariantID: "non-existent",
			},
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Variant not found"}`,
			wantErr:     true,
			errContains: "Variant not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/variants/%s",
					tt.request.ProductID, tt.request.VariantID)
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

			resp, err := UpdateProductVariant(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProductVariant() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUpdateProductImage(t *testing.T) {
	tests := []struct {
		name        string
		request     UpdateProductImageRequest
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful update",
			request: UpdateProductImageRequest{
				ProductID: "product-123",
				ImageID:   "image-123",
				AltText:   "New caption",
			},
			mockStatus: http.StatusOK,
			mockResp: `{
				"id": "image-123",
				"caption": "New caption",
				"url": "https://example.com/image.jpg"
			}`,
			wantErr: false,
		},
		{
			name: "image not found",
			request: UpdateProductImageRequest{
				ProductID: "product-123",
				ImageID:   "non-existent",
				AltText:   "Test caption",
			},
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Image not found"}`,
			wantErr:     true,
			errContains: "Image not found",
		},
		{
			name: "invalid request",
			request: UpdateProductImageRequest{
				ProductID: "",
				ImageID:   "",
				AltText:   "",
			},
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"ERROR","message":"Invalid request parameters"}`,
			wantErr:     true,
			errContains: "Invalid request parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/images/%s",
					tt.request.ProductID, tt.request.ImageID)
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

			resp, err := UpdateProductImage(context.Background(), config, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProductImage() error = %v, wantErr %v", err, tt.wantErr)
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

			if !tt.wantErr && resp != nil {
				if resp.ID != "image-123" {
					t.Errorf("expected image ID 'image-123', got %s", resp.ID)
				}
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful deletion",
			productID:  "product-123",
			mockStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:        "product not found",
			productID:   "non-existent",
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Product not found"}`,
			wantErr:     true,
			errContains: "Product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}
				expectedPath := "/1.0/commerce/products/" + tt.productID
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
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

			status, err := DeleteProduct(context.Background(), config, tt.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteProduct() error = %v, wantErr %v", err, tt.wantErr)
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

func TestDeleteProductVariant(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		variantID   string
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful deletion",
			productID:  "product-123",
			variantID:  "variant-123",
			mockStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:        "variant not found",
			productID:   "product-123",
			variantID:   "non-existent",
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Variant not found"}`,
			wantErr:     true,
			errContains: "Variant not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/variants/%s", tt.productID, tt.variantID)
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
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

			status, err := DeleteProductVariant(context.Background(), config, tt.productID, tt.variantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteProductVariant() error = %v, wantErr %v", err, tt.wantErr)
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

func TestDeleteProductImage(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		imageID     string
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful deletion",
			productID:  "product-123",
			imageID:    "image-123",
			mockStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:        "image not found",
			productID:   "product-123",
			imageID:     "non-existent",
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"ERROR","message":"Image not found"}`,
			wantErr:     true,
			errContains: "Image not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/1.0/commerce/products/%s/images/%s", tt.productID, tt.imageID)
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("expected path to end with %s, got %s", expectedPath, r.URL.Path)
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

			status, err := DeleteProductImage(context.Background(), config, tt.productID, tt.imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteProductImage() error = %v, wantErr %v", err, tt.wantErr)
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
