package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/j-low/gocommerce/common"
)

func RetrieveAllInventory(ctx context.Context, config *common.Config, params common.QueryParams) (*RetrieveAllInventoryResponse, error) {
	if err := common.ValidateQueryParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/inventory", InventoryAPIVersion)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	query := u.Query()
	if params.Cursor != "" {
		query.Set("cursor", params.Cursor)
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
		return nil, common.ParseErrorResponse("RetrieveAllInventory", u.String(), body, resp.StatusCode)
	}

	var response RetrieveAllInventoryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveSpecificInventory(ctx context.Context, config *common.Config, inventoryIDs []string) (*RetrieveSpecificInventoryResponse, error) {
	if len(inventoryIDs) == 0 {
		return nil, fmt.Errorf("no inventory IDs provided")
	}
	if len(inventoryIDs) > 50 {
		return nil, fmt.Errorf("cannot retrieve more than 50 inventory IDs")
	}

	idsPath := strings.Join(inventoryIDs, ",")
	baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/inventory", InventoryAPIVersion)
	endpoint := fmt.Sprintf("%s/%s", baseURL, idsPath)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve specific inventory: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("RetrieveSpecificInventory", u.String(), body, resp.StatusCode)
	}

	var response RetrieveSpecificInventoryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func AdjustStockQuantities(ctx context.Context, config *common.Config, request AdjustStockQuantitiesRequest) (int, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/inventory/adjustments", InventoryAPIVersion)

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

	if config.IdempotencyKey != nil {
		req.Header.Set("Idempotency-Key", config.IdempotencyKey.String())
	}

	resp, err := config.Client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to adjust stock quantities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return http.StatusNoContent, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to read response body: %w", err)
	}
	return resp.StatusCode, common.ParseErrorResponse("AdjustStockQuantities", url, body, resp.StatusCode)
}
