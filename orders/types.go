package orders

const (
	OrdersAPIVersion = "1.0"
)

type CreateOrderRequest struct {
  Items []OrderItem `json:"items"`
}

type CreateOrderResponse struct {
  OrderID string `json:"orderId"`
}

type FulfillOrderRequest struct {
  TrackingNumber string `json:"trackingNumber,omitempty"`
  Carrier        string `json:"carrier,omitempty"`
}

type RetrieveAllOrdersResponse struct {
  Orders []Order `json:"orders"`
}

type RetrieveSingleOrderResponse struct {
  Order Order `json:"order"`
}

type Order struct {
  ID     string       `json:"id"`
  Items  []OrderItem  `json:"items"`
  Total  OrderTotal   `json:"total"`
  Status string       `json:"status"`
}

type OrderItem struct {
  ProductID string `json:"productId"`
  VariantID string `json:"variantId"`
  Quantity  int    `json:"quantity"`
}

type OrderTotal struct {
  Currency string `json:"currency"`
  Amount   string `json:"amount"`
}
