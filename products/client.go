package products

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetProducts(ctx context.Context, apiKey string) ([]Product, error) {
	fmt.Printf("Getting products from Commerce API...\n")

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products", ProductsAPIVersion)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("User-Agent", "gocommerce/client")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var products GetProductsResponse
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return products.Products, nil
}

func CreateProduct(ctx context.Context, request CreateProductRequest, apiKey string) (*Product, error) {
	fmt.Printf("Creating product in Commerce API...\n")

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products", ProductsAPIVersion)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("User-Agent", "gocommerce/client")
	req.Header.Set("Content-Type", "application/json")

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(request)
	req.Body = io.NopCloser(&requestBody)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var product Product
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &product, nil
}

func UpdateProduct(ctx context.Context, productID string, request UpdateProductRequest, apiKey string) (*UpdateProductResponse, error) {
	fmt.Printf("Updating product %s in Commerce API...\n", productID)
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s", ProductsAPIVersion, productID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("User-Agent", "gocommerce/client")
	req.Header.Set("Content-Type", "application/json")

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(request)
	req.Body = io.NopCloser(&requestBody)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var updatedProduct UpdateProductResponse
	if err := json.Unmarshal(body, &updatedProduct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &updatedProduct, nil
}

func DeleteProduct(ctx context.Context, productID string, apiKey string) error {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s", ProductsAPIVersion, productID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("User-Agent", "gocommerce/client")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func CreateVariant(ctx context.Context, request CreateVariantRequest, apiKey string) (*CreateVariantResponse, error) {
	return nil, nil
}

func UpdateVariant(ctx context.Context, request UpdateVariantRequest, apiKey string) (*UpdateVariantResponse, error) {
	fmt.Printf("Updating product variant %s in Commerce API...\n", request.SKU)

	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s", ProductsAPIVersion, request.ProductID, request.VariantID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("User-Agent", "gocommerce/client")
	req.Header.Set("Content-Type", "application/json")

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(request)
	req.Body = io.NopCloser(&requestBody)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product variant: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var updatedVariant UpdateVariantResponse
	if err := json.Unmarshal(body, &updatedVariant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &updatedVariant, nil
}

func SetVariantImage(ctx context.Context, request SetVariantImageRequest, apiKey string) error {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/products/%s/variants/%s/image", ProductsAPIVersion, request.ProductID, request.VariantID)

    reqBody, err := json.Marshal(request)
    if err != nil {
        return fmt.Errorf("failed to marshal request body: %v", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Authorization", "Bearer "+ apiKey)
    req.Header.Set("User-Agent", "gocommerce/client")
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to assign image to variant: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("unexpected response code: %d, body: %s", resp.StatusCode, string(body))
    }

    return nil
}

func GetStorePages(ctx context.Context, request GetStorePagesRequest, apiKey string) (*GetStorePagesResponse, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/store_pages", ProductsAPIVersion)
	if request.Cursor != "" {
		url = fmt.Sprintf("%s?cursor=%s", url, request.Cursor)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+ apiKey)
	req.Header.Set("User-Agent", "gocommerce/client")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store pages: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response GetStorePagesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &response, nil
}