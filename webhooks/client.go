package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/j-low/gocommerce/common"
)

func CreateWebhookSubscription(ctx context.Context, config *common.Config, request WebhookSubscriptionRequest) (*WebhookSubscription, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	baseURL, err := common.BuildBaseURL(config, WebhooksAPIVersion, "webhook_subscriptions")
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	if len(request.Topics) == 0 {
		return nil, fmt.Errorf("topics cannot be empty")
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
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
		return nil, common.ParseErrorResponse("CreateWebhookSubscription", baseURL, body, resp.StatusCode)
	}

	var response WebhookSubscription
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func UpdateWebhookSubscription(ctx context.Context, config *common.Config, subscriptionID string, request WebhookSubscriptionRequest) (*WebhookSubscription, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	baseURL, err := common.BuildBaseURL(config, WebhooksAPIVersion, fmt.Sprintf("webhook_subscriptions/%s", subscriptionID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
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
		return nil, common.ParseErrorResponse("UpdateWebhookSubscription", baseURL, body, resp.StatusCode)
	}

	var response WebhookSubscription
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveAllWebhookSubscriptions(ctx context.Context, config *common.Config) (*RetrieveAllWebhookSubscriptionsResponse, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	baseURL, err := common.BuildBaseURL(config, WebhooksAPIVersion, "webhook_subscriptions")
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
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
		return nil, common.ParseErrorResponse("RetrieveAllWebhookSubscriptions", baseURL, body, resp.StatusCode)
	}

	var response RetrieveAllWebhookSubscriptionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RetrieveSpecificWebhookSubscription(ctx context.Context, config *common.Config, subscriptionID string) (*WebhookSubscription, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	baseURL, err := common.BuildBaseURL(config, WebhooksAPIVersion, fmt.Sprintf("webhook_subscriptions/%s", subscriptionID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
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
		return nil, common.ParseErrorResponse("RetrieveSpecificWebhookSubscription", baseURL, body, resp.StatusCode)
	}

	var response WebhookSubscription
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func DeleteWebhookSubscription(ctx context.Context, config *common.Config, subscriptionID string) (int, error) {
	if config.AccessToken == "" {
		return http.StatusBadRequest, fmt.Errorf("access token is required")
	}

	if subscriptionID == "" {
		return http.StatusBadRequest, fmt.Errorf("subscriptionID cannot be empty")
	}

	baseURL, err := common.BuildBaseURL(config, WebhooksAPIVersion, fmt.Sprintf("webhook_subscriptions/%s", subscriptionID))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, baseURL, nil)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	req.Header.Set("User-Agent", common.SetUserAgent(config.UserAgent))

	resp, err := config.Client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to delete webhook subscription: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to read response body: %w", readErr)
		}
		return resp.StatusCode, common.ParseErrorResponse("DeleteWebhookSubscription", baseURL, body, resp.StatusCode)
	}

	return resp.StatusCode, nil
}

func SendTestNotification(ctx context.Context, config *common.Config, subscriptionID string, request SendTestNotificationRequest) (*SendTestNotificationResponse, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	baseURL, err := common.BuildBaseURL(config, WebhooksAPIVersion, fmt.Sprintf("webhook_subscriptions/%s/actions/sendTestNotification", subscriptionID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	if request.Topic == "" {
		return nil, fmt.Errorf("topic is required")
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
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
		return nil, common.ParseErrorResponse("SendTestNotification", baseURL, body, resp.StatusCode)
	}

	var response SendTestNotificationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}

func RotateSubscriptionSecret(ctx context.Context, config *common.Config, subscriptionID string) (*RotateSubscriptionSecretResponse, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	baseURL, err := common.BuildBaseURL(config, WebhooksAPIVersion, fmt.Sprintf("webhook_subscriptions/%s/actions/rotateSecret", subscriptionID))
	if err != nil {
		return nil, fmt.Errorf("failed to build base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
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
		return nil, common.ParseErrorResponse("RotateSubscriptionSecret", baseURL, body, resp.StatusCode)
	}

	var response RotateSubscriptionSecretResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &response, nil
}
