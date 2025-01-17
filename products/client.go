package products

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/j-low/gocommerce/common"
)

func CreateProduct(ctx context.Context, config *common.Config, request CreateProductRequest) (*Product, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products", ProductsAPIVersion)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, common.ParseErrorResponse("CreateProduct", url, body, resp.StatusCode)
	}

	var product Product
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &product, nil
}

func CreateProductVariant(ctx context.Context, config *common.Config, request CreateProductVariantRequest) (*CreateProductVariantResponse, error) {
	if request.ProductID == "" {
		return nil, fmt.Errorf("productID is required")
	}

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants", ProductsAPIVersion, request.ProductID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product variant: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, common.ParseErrorResponse("CreateProductVariant", url, body, resp.StatusCode)
	}

	var createdVariant CreateProductVariantResponse
	if err := json.Unmarshal(body, &createdVariant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &createdVariant, nil
}

func UploadProductImage(ctx context.Context, config *common.Config, productID, filePath string) (*UploadProductImageResponse, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images", ProductsAPIVersion, productID)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload product image: %w", err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %w", readErr)
	}
	if resp.StatusCode != http.StatusAccepted {
		return nil, common.ParseErrorResponse("UploadProductImage", url, body, resp.StatusCode)
	}

	var response UploadProductImageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveAllStorePages(ctx context.Context, config *common.Config, params common.QueryParams) (*RetrieveAllStorePagesResponse, error) {
	if err := common.ValidateQueryParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}
	baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/store_pages", ProductsAPIVersion)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	if params.Cursor != "" {
		query := u.Query()
		query.Set("cursor", params.Cursor)
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store pages: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("RetrieveAllStorePages", u.String(), body, resp.StatusCode)
	}

	var response RetrieveAllStorePagesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveAllProducts(ctx context.Context, config *common.Config, params common.QueryParams) (*RetrieveAllProductsResponse, error) {
	if err := common.ValidateQueryParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products", ProductsAPIVersion)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	query := u.Query()
	if params.Cursor != "" {
		query.Set("cursor", params.Cursor)
	}
	if params.ModifiedAfter != "" {
		query.Set("modifiedAfter", params.ModifiedAfter)
	}
	if params.ModifiedBefore != "" {
		query.Set("modifiedBefore", params.ModifiedBefore)
	}
	if params.Type != "" {
		query.Set("type", params.Type)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("RetrieveAllProducts", u.String(), body, resp.StatusCode)
	}

	var response RetrieveAllProductsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveSpecificProducts(ctx context.Context, config *common.Config, productIDs []string) (*RetrieveSpecificProductsResponse, error) {
	if len(productIDs) == 0 {
		return nil, fmt.Errorf("at least one product ID is required")
	}
	if len(productIDs) > 50 {
		return nil, fmt.Errorf("cannot retrieve more than 50 products at once")
	}

	joinedIDs := strings.Join(productIDs, ",")

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s", ProductsAPIVersion, joinedIDs)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve specific products: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("RetrieveSpecificProducts", url, body, resp.StatusCode)
	}

	var response RetrieveSpecificProductsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func GetProductImageUploadStatus(ctx context.Context, config *common.Config, productID, imageID string) (*GetProductImageUploadStatusResponse, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images/%s/status", ProductsAPIVersion, productID, imageID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product image upload status: %w", err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %w", readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("GetProductImageUploadStatus", url, body, resp.StatusCode)
	}

	var statusResponse GetProductImageUploadStatusResponse
	if err := json.Unmarshal(body, &statusResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &statusResponse, nil
}

func AssignProductImageToVariant(ctx context.Context, config *common.Config, request AssignProductImageToVariantRequest) (int, error) {
	if ctx == nil {
		return http.StatusBadRequest, fmt.Errorf("context cannot be nil")
	}

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s/image", ProductsAPIVersion, request.ProductID, request.VariantID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to assign image to variant: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to read response body: %w", err)
		}
		return resp.StatusCode, common.ParseErrorResponse("AssignProductImageToVariant", url, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func ReorderProductImage(ctx context.Context, config *common.Config, request ReorderProductImageRequest) (int, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images/%s/order", ProductsAPIVersion, request.ProductID, request.ImageID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to reorder product image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to read response body: %w", readErr)
		}
		return resp.StatusCode, common.ParseErrorResponse("ReorderProductImage", url, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func UpdateProduct(ctx context.Context, config *common.Config, productID string, request UpdateProductRequest) (*UpdateProductResponse, error) {
	if productID == "" {
		return nil, fmt.Errorf("productID is required")
	}

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s", ProductsAPIVersion, productID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("UpdateProduct", url, body, resp.StatusCode)
	}

	var updatedProduct UpdateProductResponse
	if err := json.Unmarshal(body, &updatedProduct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &updatedProduct, nil
}

func UpdateProductVariant(ctx context.Context, config *common.Config, request UpdateProductVariantRequest) (*UpdateProductVariantResponse, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s", ProductsAPIVersion, request.ProductID, request.VariantID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product variant: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("UpdateProductVariant", url, body, resp.StatusCode)
	}

	var updatedVariant UpdateProductVariantResponse
	if err := json.Unmarshal(body, &updatedVariant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &updatedVariant, nil
}

func UpdateProductImage(ctx context.Context, config *common.Config, request UpdateProductImageRequest) (*UpdateProductImageResponse, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images/%s", ProductsAPIVersion, request.ProductID, request.ImageID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product image: %w", err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %w", readErr)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("UpdateProductImage", url, body, resp.StatusCode)
	}

	var updatedImage UpdateProductImageResponse
	if err := json.Unmarshal(body, &updatedImage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &updatedImage, nil
}

func DeleteProduct(ctx context.Context, config *common.Config, productID string) (int, error) {
	if productID == "" {
		return http.StatusBadRequest, fmt.Errorf("productID is required")
	}

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s", ProductsAPIVersion, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to delete product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to read response body: %w", err)
		}
		return resp.StatusCode, common.ParseErrorResponse("DeleteProduct", url, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func DeleteProductVariant(ctx context.Context, config *common.Config, productID, variantID string) (int, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s", ProductsAPIVersion, productID, variantID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to delete product variant: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to read response body: %w", readErr)
		}
		return resp.StatusCode, common.ParseErrorResponse("DeleteProductVariant", url, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func DeleteProductImage(ctx context.Context, config *common.Config, productID, imageID string) (int, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images/%s", ProductsAPIVersion, productID, imageID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to delete product image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to read response body: %w", readErr)
		}
		return resp.StatusCode, common.ParseErrorResponse("DeleteProductImage", url, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}
