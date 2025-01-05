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

	"github.com/NuvoCodeTechnologies/gocommerce/common"
)

func CreateProduct(ctx context.Context, config *common.Config, request CreateProductRequest) (*Product, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products", ProductsAPIVersion)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
		return nil, common.ParseErrorResponse(body, resp.StatusCode)
	}

	var product Product
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &product, nil
}

func CreateVariant(ctx context.Context, config *common.Config, request CreateVariantRequest) (*CreateVariantResponse, error) {
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

  req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var createdVariant CreateVariantResponse
  if err := json.Unmarshal(body, &createdVariant); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &createdVariant, nil
}

func UploadProductImage(ctx context.Context, config *common.Config, request UploadProductImageRequest) (*UploadProductImageResponse, error) {
  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images", ProductsAPIVersion, request.ProductID)

  file, err := os.Open(request.FilePath)
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

  req, err := http.NewRequestWithContext(ctx, "POST", url, &requestBody)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response UploadProductImageResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func GetStorePages(ctx context.Context, config *common.Config, request GetStorePagesRequest) (*GetStorePagesResponse, error) {
	baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/store_pages", ProductsAPIVersion)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	if request.Cursor != "" {
		query := u.Query()
		query.Set("cursor", request.Cursor)
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
		return nil, common.ParseErrorResponse(body, resp.StatusCode)
	}

	var response GetStorePagesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func GetAllProducts(ctx context.Context, config *common.Config, request GetAllProductsRequest) (*GetAllProductsResponse, error) {
	baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products", ProductsAPIVersion)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	
	queryParams := url.Values{}
	if request.Cursor != "" {
		queryParams.Add("cursor", request.Cursor)
	}
	if request.ModifiedAfter != "" {
		queryParams.Add("modifiedAfter", request.ModifiedAfter)
	}
	if request.ModifiedBefore != "" {
		queryParams.Add("modifiedBefore", request.ModifiedBefore)
	}
	if request.Type != "" {
		queryParams.Add("type", request.Type)
	}

	query := u.Query()
	for key, value := range queryParams {
		query.Add(key, value[0])
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
		return nil, common.ParseErrorResponse(body, resp.StatusCode)
	}

	var response GetAllProductsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	
	return &response, nil
}

func GetSpecificProducts(ctx context.Context, config *common.Config, request GetSpecificProductsRequest) (*GetSpecificProductsResponse, error) {
  if len(request.ProductIDs) == 0 {
    return nil, fmt.Errorf("at least one product ID is required")
  }

  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products", ProductsAPIVersion)

  reqBody, err := json.Marshal(request)
  if err != nil {
    return nil, fmt.Errorf("failed to marshal request body: %w", err)
  }

  req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer " + config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
  req.Header.Set("Content-Type", "application/json")

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
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response GetSpecificProductsResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func GetProductImageUploadStatus(ctx context.Context, config *common.Config, productID, imageID string) (*GetProductImageUploadStatusResponse, error) {
  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images/%s/status", ProductsAPIVersion, productID, imageID)

  req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var statusResponse GetProductImageUploadStatusResponse
  if err := json.Unmarshal(body, &statusResponse); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &statusResponse, nil
}

func AssignProductImageToVariant(ctx context.Context, config *common.Config, request AssignProductImageToVariantRequest) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s/image", ProductsAPIVersion, request.ProductID, request.VariantID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to assign image to variant: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	return common.ParseErrorResponse(body, resp.StatusCode)
	}

	return nil
}

func ReorderProductImage(ctx context.Context, config *common.Config, request ReorderProductImageRequest) error {
  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images/%s/order", ProductsAPIVersion, request.ProductID, request.ImageID)

  reqBody, err := json.Marshal(request)
  if err != nil {
    return fmt.Errorf("failed to marshal request body: %w", err)
  }

  req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
  if err != nil {
    return fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer " + config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
  req.Header.Set("Content-Type", "application/json")

  resp, err := config.Client.Do(req)
  if err != nil {
    return fmt.Errorf("failed to reorder product image: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusNoContent {
    body, readErr := io.ReadAll(resp.Body)
    if readErr != nil {
      return fmt.Errorf("failed to read response body: %w", readErr)
    }
    return common.ParseErrorResponse(body, resp.StatusCode)
  }

  return nil
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

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
		return nil, common.ParseErrorResponse(body, resp.StatusCode)
	}

	var updatedProduct UpdateProductResponse
	if err := json.Unmarshal(body, &updatedProduct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &updatedProduct, nil
}

func UpdateVariant(ctx context.Context, config *common.Config, request UpdateVariantRequest) (*UpdateVariantResponse, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s", ProductsAPIVersion, request.ProductID, request.VariantID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
		return nil, common.ParseErrorResponse(body, resp.StatusCode)
	}

	var updatedVariant UpdateVariantResponse
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

  req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var updatedImage UpdateProductImageResponse
  if err := json.Unmarshal(body, &updatedImage); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &updatedImage, nil
}

func DeleteProduct(ctx context.Context, config *common.Config, productID string) error {
	if productID == "" {
		return fmt.Errorf("productID is required")
	}

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s", ProductsAPIVersion, productID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return common.ParseErrorResponse(body, resp.StatusCode)
	}

	return nil
}

func DeleteVariant(ctx context.Context, config *common.Config, productID, variantID string) error {
  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s", ProductsAPIVersion, productID, variantID)

  req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
  if err != nil {
    return fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Authorization", "Bearer " + config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

  resp, err := config.Client.Do(req)
  if err != nil {
    return fmt.Errorf("failed to delete product variant: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusNoContent {
    body, readErr := io.ReadAll(resp.Body)
    if readErr != nil {
      return fmt.Errorf("failed to read response body: %w", readErr)
    }
    return common.ParseErrorResponse(body, resp.StatusCode)
  }

  return nil
}

func DeleteProductImage(ctx context.Context, config *common.Config, productID, imageID string) error {
  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/images/%s", ProductsAPIVersion, productID, imageID)

  req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
  if err != nil {
    return fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Authorization", "Bearer " + config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

  resp, err := config.Client.Do(req)
  if err != nil {
    return fmt.Errorf("failed to delete product image: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusNoContent {
    body, readErr := io.ReadAll(resp.Body)
    if readErr != nil {
      return fmt.Errorf("failed to read response body: %w", readErr)
    }
    return common.ParseErrorResponse(body, resp.StatusCode)
  }

  return nil
}
