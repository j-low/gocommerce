package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/NuvoCodeTechnologies/gocommerce/common"
)

func RetrieveAllInventory(ctx context.Context, config *common.Config, queryParams common.QueryParams) (*RetrieveAllInventoryResponse, error) {
  baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/inventory", InventoryAPIVersion)
  u, err := url.Parse(baseURL)
  if err != nil {
    return nil, fmt.Errorf("failed to parse base URL: %w", err)
  }

  query := url.Values{}
  if queryParams.Cursor != "" {
    query.Add("cursor", queryParams.Cursor)
  }

  u.RawQuery = query.Encode()

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Authorization", "Bearer " + config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

  resp, err := config.Client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to retrieve all inventory: %w", err)
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response RetrieveAllInventoryResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func RetrieveSpecificInventory(ctx context.Context, config *common.Config, request RetrieveSpecificInventoryRequest) (*RetrieveSpecificInventoryResponse, error) {
  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/inventory/bulk", InventoryAPIVersion)

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
    return nil, fmt.Errorf("failed to retrieve specific inventory: %w", err)
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response RetrieveSpecificInventoryResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func AdjustStockQuantities(ctx context.Context, config *common.Config, request AdjustStockQuantitiesRequest) (*AdjustStockQuantitiesResponse, error) {
  url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/inventory/adjustments", InventoryAPIVersion)

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
    return nil, fmt.Errorf("failed to adjust stock quantities: %w", err)
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response AdjustStockQuantitiesResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}
