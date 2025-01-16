package orders

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

func CreateOrder(ctx context.Context, config *common.Config, request CreateOrderRequest) (*Order, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/orders", OrdersAPIVersion)

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

	if config.IdempotencyKey != nil {
		req.Header.Set("Idempotency-Key", config.IdempotencyKey.String())
	}

	resp, err := config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, common.ParseErrorResponse(body, resp.StatusCode)
	}

	var response Order
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func FulfillOrder(ctx context.Context, config *common.Config, orderID string, request FulfillOrderRequest) error {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/orders/%s/fulfillments", OrdersAPIVersion, orderID)

	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + config.APIKey)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fulfill order: %w", err)
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

func RetrieveAllOrders(ctx context.Context, config *common.Config, queryParams common.QueryParams) (*RetrieveAllOrdersResponse, error) {
	if queryParams.Cursor != "" {
		if queryParams.ModifiedAfter != "" || queryParams.ModifiedBefore != "" {
			return nil, fmt.Errorf("cannot use both cursor and modifiedAfter/modifiedBefore")
		}
	} else {
		if (queryParams.ModifiedAfter != "" && queryParams.ModifiedBefore == "") ||
			(queryParams.ModifiedAfter == "" && queryParams.ModifiedBefore != "") {
			return nil, fmt.Errorf("must specify both modifiedAfter and modifiedBefore, or neither")
		}
	}

	baseURL := fmt.Sprintf("https://api.squarespace.com/%s/commerce/orders", OrdersAPIVersion)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	query := u.Query()
	if queryParams.Cursor != "" {
		query.Set("cursor", queryParams.Cursor)
	}
	if queryParams.ModifiedAfter != "" {
		query.Set("modifiedAfter", queryParams.ModifiedAfter)
	}
	if queryParams.ModifiedBefore != "" {
		query.Set("modifiedBefore", queryParams.ModifiedBefore)
	}
	if queryParams.Status != "" {
		query.Set("fulfillmentStatus", queryParams.Status)
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
		return nil, fmt.Errorf("failed to retrieve all orders: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse(body, resp.StatusCode)
	}

	var response RetrieveAllOrdersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveSpecificOrder(ctx context.Context, config *common.Config, orderID string) (*Order, error) {
	url := fmt.Sprintf("https://api.squarespace.com/%s/commerce/orders/%s", OrdersAPIVersion, orderID)

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer " + config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

  resp, err := config.Client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to retrieve order: %w", err)
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response Order
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}
