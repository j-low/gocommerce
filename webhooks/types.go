package webhooks

const (
	WebhooksAPIVersion = "1.0"
)

type WebhookSubscriptionRequest struct {
	EndpointURL string   `json:"endpointUrl"`
	Topics      []string `json:"topics"`
}

type RetrieveAllWebhookSubscriptionsResponse struct {
	WebhookSubscriptions []WebhookSubscription `json:"webhookSubscriptions"`
}

type SendTestNotificationRequest struct {
	Topic string `json:"topic"`
}

type SendTestNotificationResponse struct {
	StatusCode int `json:"statusCode"`
}

type RotateSubscriptionSecretResponse struct {
	Secret string `json:"secret"`
}

type WebhookSubscription struct {
	ID          string   `json:"id"`
	EndpointURL string   `json:"endpointUrl"`
	Topics      []string `json:"topics"`
	Secret      string   `json:"secret"`
	CreatedOn   string   `json:"createdOn"`
	UpdatedOn   string   `json:"updatedOn"`
}
