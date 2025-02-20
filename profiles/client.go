package profiles

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

func RetrieveAllProfiles(ctx context.Context, config *common.Config, params common.QueryParams) (*RetrieveAllProfilesResponse, error) {
	if err := common.ValidateQueryParams(params); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	baseURL, err := common.BuildBaseURL(config, ProfilesAPIVersion, "profiles")
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
	if params.Filter != "" {
		query.Set("filter", params.Filter)
	}
	if params.SortDirection != "" {
		query.Set("sortDirection", params.SortDirection)
	}
	if params.SortField != "" {
		query.Set("sortField", params.SortField)
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
		return nil, fmt.Errorf("failed to retrieve all profiles: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("RetrieveAllProfiles", u.String(), body, resp.StatusCode)
	}

	var response RetrieveAllProfilesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveSpecificProfiles(ctx context.Context, config *common.Config, profileIDs []string) (*RetrieveSpecificProfilesResponse, error) {
	if len(profileIDs) == 0 {
		return nil, fmt.Errorf("profileIDs cannot be empty")
	}

	joinedIDs := strings.Join(profileIDs, ",")

	baseURL, err := common.BuildBaseURL(config, ProfilesAPIVersion, fmt.Sprintf("profiles/%s", joinedIDs))
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
		return nil, fmt.Errorf("failed to retrieve specific profiles: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ParseErrorResponse("RetrieveSpecificProfiles", u.String(), body, resp.StatusCode)
	}

	var response RetrieveSpecificProfilesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}
