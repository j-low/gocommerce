package transactions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/j-low/gocommerce/common"
)

func RetrieveAllTransactions(ctx context.Context, config *common.Config, params common.QueryParams) (*RetrieveAllTransactionsResponse, error) {
	if err := common.ValidateQueryParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	baseURL, err := common.BuildBaseURL(config, TransactionsAPIVersion, "commerce/transactions")
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
		return nil, common.ParseErrorResponse("RetrieveAllTransactions", u.String(), body, resp.StatusCode)
	}

	var response RetrieveAllTransactionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveSpecificTransactions(ctx context.Context, config *common.Config, transactionIDs []string) (*RetrieveSpecificTransactionsResponse, error) {
	if len(transactionIDs) == 0 {
		return nil, fmt.Errorf("transactionIDs cannot be empty")
	}
	if len(transactionIDs) > 50 {
		return nil, fmt.Errorf("transactionIDs cannot exceed 50 IDs")
	}

	ids := url.PathEscape(strings.Join(transactionIDs, ","))
	baseURL, err := common.BuildBaseURL(config, TransactionsAPIVersion, fmt.Sprintf("commerce/transactions/%s", ids))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

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
		return nil, common.ParseErrorResponse("RetrieveSpecificTransactions", baseURL, body, resp.StatusCode)
	}

	var response RetrieveSpecificTransactionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}
