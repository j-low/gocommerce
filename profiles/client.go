package profiles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/NuvoCodeTechnologies/gocommerce/common"
)

func RetrieveAllProfiles(ctx context.Context, config *common.Config, params common.QueryParams) (*RetrieveAllProfilesResponse, error) {
  baseURL := fmt.Sprintf("https://api.squarespace.com/%s/profiles", ProfilesAPIVersion)
  u, err := url.Parse(baseURL)
  if err != nil {
    return nil, fmt.Errorf("failed to parse base URL: %w", err)
  }

  queryParams := url.Values{}
  if params.Cursor != "" {
    queryParams.Add("cursor", params.Cursor)
  }
  if params.Filter != "" {
    queryParams.Add("filter", params.Filter)
  }
  if params.SortDirection != "" {
    queryParams.Add("sortDirection", params.SortDirection)
  }
  if params.SortField != "" {
    queryParams.Add("sortField", params.SortField)
  }

  u.RawQuery = queryParams.Encode()

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
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

  baseURL := fmt.Sprintf("https://api.squarespace.com/%s/profiles", ProfilesAPIVersion)
  u, err := url.Parse(baseURL)
  if err != nil {
    return nil, fmt.Errorf("failed to parse base URL: %w", err)
  }

  queryParams := url.Values{}
  queryParams.Add("ids", url.QueryEscape(fmt.Sprintf("%s", profileIDs)))

  u.RawQuery = queryParams.Encode()

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Authorization", "Bearer " + config.APIKey)
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
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response RetrieveSpecificProfilesResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}
