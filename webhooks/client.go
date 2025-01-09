package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NuvoCodeTechnologies/gocommerce/common"
)

func CreateWebhookSubscription(ctx context.Context, config *common.Config, request WebhookSubscriptionRequest) (*WebhookSubscription, error) {
  url := fmt.Sprintf("https://api.squarespace.com/%s/webhook_subscriptions", WebhooksAPIVersion)

  if len(request.Topics) == 0 {
    return nil, fmt.Errorf("topics cannot be empty")
  }

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
    return nil, fmt.Errorf("failed to create webhook subscription: %w", err)
  }
  defer resp.Body.Close()

  body, readErr := io.ReadAll(resp.Body)
  if readErr != nil {
    return nil, fmt.Errorf("failed to read response body: %w", readErr)
  }

  if resp.StatusCode != http.StatusCreated {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response WebhookSubscription
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func UpdateWebhookSubscription(ctx context.Context, config *common.Config, subscriptionID string, request WebhookSubscriptionRequest) (*WebhookSubscription, error) {
  url := fmt.Sprintf("https://api.squarespace.com/%s/webhook_subscriptions/%s", WebhooksAPIVersion, subscriptionID)

  if subscriptionID == "" {
    return nil, fmt.Errorf("subscriptionID cannot be empty")
  }

  if request.Topics != nil && len(request.Topics) == 0 {
    return nil, fmt.Errorf("topics cannot be an empty array")
  }

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
    return nil, fmt.Errorf("failed to update webhook subscription: %w", err)
  }
  defer resp.Body.Close()

  body, readErr := io.ReadAll(resp.Body)
  if readErr != nil {
    return nil, fmt.Errorf("failed to read response body: %w", readErr)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response WebhookSubscription
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func RetrieveAllWebhookSubscriptions(ctx context.Context, config *common.Config) (*RetrieveAllWebhookSubscriptionsResponse, error) {
  url := fmt.Sprintf("https://api.squarespace.com/%s/webhook_subscriptions", WebhooksAPIVersion)

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer "+config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

  resp, err := config.Client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to retrieve webhook subscriptions: %w", err)
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response RetrieveAllWebhookSubscriptionsResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func RetrieveSpecificWebhookSubscription(ctx context.Context, config *common.Config, subscriptionID string) (*WebhookSubscription, error) {
  if subscriptionID == "" {
    return nil, fmt.Errorf("subscriptionID cannot be empty")
  }

  url := fmt.Sprintf("https://api.squarespace.com/%s/webhook_subscriptions/%s", WebhooksAPIVersion, subscriptionID)

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer "+config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

  resp, err := config.Client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to retrieve webhook subscription: %w", err)
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response WebhookSubscription
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func DeleteWebhookSubscription(ctx context.Context, config *common.Config, subscriptionID string) error {
  if subscriptionID == "" {
    return fmt.Errorf("subscriptionID cannot be empty")
  }

  url := fmt.Sprintf("https://api.squarespace.com/%s/webhook_subscriptions/%s", WebhooksAPIVersion, subscriptionID)

  req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
  if err != nil {
    return fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer "+config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

  resp, err := config.Client.Do(req)
  if err != nil {
    return fmt.Errorf("failed to delete webhook subscription: %w", err)
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

func SendTestNotification(ctx context.Context, config *common.Config, subscriptionID string, request SendTestNotificationRequest) (*SendTestNotificationResponse, error) {
  if subscriptionID == "" {
    return nil, fmt.Errorf("subscriptionID cannot be empty")
  }

  url := fmt.Sprintf("https://api.squarespace.com/%s/webhook_subscriptions/%s/actions/sendTestNotification", WebhooksAPIVersion, subscriptionID)

  if request.Topic == "" {
    return nil, fmt.Errorf("topic is required")
  }

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
    return nil, fmt.Errorf("failed to send test notification: %w", err)
  }
  defer resp.Body.Close()

  body, readErr := io.ReadAll(resp.Body)
  if readErr != nil {
    return nil, fmt.Errorf("failed to read response body: %w", readErr)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response SendTestNotificationResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}

func RotateSubscriptionSecret(ctx context.Context, config *common.Config, subscriptionID string) (*RotateSubscriptionSecretResponse, error) {
  if subscriptionID == "" {
    return nil, fmt.Errorf("subscriptionID cannot be empty")
  }

  url := fmt.Sprintf("https://api.squarespace.com/%s/webhook_subscriptions/%s/actions/rotateSecret", WebhooksAPIVersion, subscriptionID)

  req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }

  req.Header.Set("Authorization", "Bearer "+config.APIKey)
  req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))
  req.Header.Set("Content-Type", "application/json")

  resp, err := config.Client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to rotate subscription secret: %w", err)
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  if resp.StatusCode != http.StatusOK {
    return nil, common.ParseErrorResponse(body, resp.StatusCode)
  }

  var response RotateSubscriptionSecretResponse
  if err := json.Unmarshal(body, &response); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
  }

  return &response, nil
}