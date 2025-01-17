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
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, "commerce/products")
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
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
		return nil, common.ParseErrorResponse("CreateProduct", baseURL, body, resp.StatusCode)
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

	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/variants", request.ProductID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
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
		return nil, common.ParseErrorResponse("CreateProductVariant", baseURL, body, resp.StatusCode)
	}

	var createdVariant CreateProductVariantResponse
	if err := json.Unmarshal(body, &createdVariant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &createdVariant, nil
}

func UploadProductImage(ctx context.Context, config *common.Config, productID, filePath string) (*UploadProductImageResponse, error) {
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/images", productID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, &requestBody)
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
		return nil, common.ParseErrorResponse("UploadProductImage", baseURL, body, resp.StatusCode)
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

	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, "commerce/store_pages")
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

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

	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, "commerce/products")
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

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

	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s", joinedIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
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
		return nil, common.ParseErrorResponse("RetrieveSpecificProducts", baseURL, body, resp.StatusCode)
	}

	var response RetrieveSpecificProductsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func GetProductImageUploadStatus(ctx context.Context, config *common.Config, productID, imageID string) (*GetProductImageUploadStatusResponse, error) {
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/images/%s/status", productID, imageID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
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
		return nil, common.ParseErrorResponse("GetProductImageUploadStatus", baseURL, body, resp.StatusCode)
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

	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/variants/%s/image", request.ProductID, request.VariantID))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to build base URL: %w", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
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
		return resp.StatusCode, common.ParseErrorResponse("AssignProductImageToVariant", baseURL, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func ReorderProductImage(ctx context.Context, config *common.Config, request ReorderProductImageRequest) (int, error) {
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/images/%s/order", request.ProductID, request.ImageID))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to build base URL: %w", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
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
		return resp.StatusCode, common.ParseErrorResponse("ReorderProductImage", baseURL, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func UpdateProduct(ctx context.Context, config *common.Config, productID string, request UpdateProductRequest) (*UpdateProductResponse, error) {
	if productID == "" {
		return nil, fmt.Errorf("productID is required")
	}

	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s", productID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
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
		return nil, common.ParseErrorResponse("UpdateProduct", baseURL, body, resp.StatusCode)
	}

	var updatedProduct UpdateProductResponse
	if err := json.Unmarshal(body, &updatedProduct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &updatedProduct, nil
}

func UpdateProductVariant(ctx context.Context, config *common.Config, request UpdateProductVariantRequest) (*UpdateProductVariantResponse, error) {
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/variants/%s", request.ProductID, request.VariantID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
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
		return nil, common.ParseErrorResponse("UpdateProductVariant", baseURL, body, resp.StatusCode)
	}

	var updatedVariant UpdateProductVariantResponse
	if err := json.Unmarshal(body, &updatedVariant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &updatedVariant, nil
}

func UpdateProductImage(ctx context.Context, config *common.Config, request UpdateProductImageRequest) (*UpdateProductImageResponse, error) {
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/images/%s", request.ProductID, request.ImageID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
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
		return nil, common.ParseErrorResponse("UpdateProductImage", baseURL, body, resp.StatusCode)
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

	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s", productID))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL, nil)
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
		return resp.StatusCode, common.ParseErrorResponse("DeleteProduct", baseURL, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func DeleteProductVariant(ctx context.Context, config *common.Config, productID, variantID string) (int, error) {
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/variants/%s", productID, variantID))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL, nil)
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
		return resp.StatusCode, common.ParseErrorResponse("DeleteProductVariant", baseURL, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}

func DeleteProductImage(ctx context.Context, config *common.Config, productID, imageID string) (int, error) {
	baseURL, err := common.BuildBaseURL(config, ProductsAPIVersion, fmt.Sprintf("commerce/products/%s/images/%s", productID, imageID))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL, nil)
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
		return resp.StatusCode, common.ParseErrorResponse("DeleteProductImage", baseURL, body, resp.StatusCode)
	}

	return http.StatusNoContent, nil
}
