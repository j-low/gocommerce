package orders

import "github.com/NuvoCodeTechnologies/gocommerce/common"

const (
	OrdersAPIVersion = "1.0"
)

type CreateOrderRequest struct {
  ChannelName          string               `json:"channelName"`
  ExternalOrderReference string             `json:"externalOrderReference"`
  CustomerEmail        *string              `json:"customerEmail,omitempty"`
  BillingAddress       *Address             `json:"billingAddress,omitempty"`
  ShippingAddress      *Address             `json:"shippingAddress,omitempty"`
  InventoryBehavior    *string              `json:"inventoryBehavior,omitempty"`
  LineItems            []LineItem           `json:"lineItems"`
  ShippingLines        []ShippingLine       `json:"shippingLines,omitempty"`
  DiscountLines        []DiscountLine       `json:"discountLines,omitempty"`
  PriceTaxInterpretation string             `json:"priceTaxInterpretation"`
  Subtotal             *common.Amount       `json:"subtotal,omitempty"`
  ShippingTotal        *common.Amount       `json:"shippingTotal,omitempty"`
  DiscountTotal        *common.Amount       `json:"discountTotal,omitempty"`
  TaxTotal             *common.Amount       `json:"taxTotal,omitempty"`
  GrandTotal           common.Amount        `json:"grandTotal"`
  FulfillmentStatus    *string              `json:"fulfillmentStatus,omitempty"`
  ShopperFulfillmentNotificationBehavior *string `json:"shopperFulfillmentNotificationBehavior,omitempty"`
  FulfilledOn          *string              `json:"fulfilledOn,omitempty"`
  Fulfillments         []Fulfillment        `json:"fulfillments"`
  CreatedOn            string               `json:"createdOn"`
}

type CreateOrderResponse struct {
  OrderID string `json:"orderId"`
}

type FulfillOrderRequest struct {
  TrackingNumber string `json:"trackingNumber,omitempty"`
  Carrier        string `json:"carrier,omitempty"`
}

type RetrieveAllOrdersResponse struct {
  Result     []Order    `json:"result"`
  Pagination common.Pagination `json:"pagination"`
}

type RetrieveSingleOrderResponse struct {
  Order Order `json:"order"`
}

type Order struct {
  ID                    string          `json:"id"`
  OrderNumber           string          `json:"orderNumber"`
  CreatedOn             string          `json:"createdOn"`
  ModifiedOn            string          `json:"modifiedOn"`
  Channel               string          `json:"channel"`
  TestMode              bool            `json:"testmode"`
  CustomerEmail         string          `json:"customerEmail"`
  BillingAddress        Address         `json:"billingAddress"`
  ShippingAddress       Address         `json:"shippingAddress"`
  FulfillmentStatus     string          `json:"fulfillmentStatus"`
  LineItems             []LineItem      `json:"lineItems"`
  InternalNotes         []Note          `json:"internalNotes"`
  ShippingLines         []ShippingLine  `json:"shippingLines"`
  DiscountLines         []DiscountLine  `json:"discountLines"`
  FormSubmission        []FormSubmission `json:"formSubmission"`
  Fulfillments          []Fulfillment   `json:"fulfillments"`
  Subtotal              common.Amount          `json:"subtotal"`
  ShippingTotal         common.Amount          `json:"shippingTotal"`
  DiscountTotal         common.Amount          `json:"discountTotal"`
  TaxTotal              common.Amount          `json:"taxTotal"`
  RefundedTotal         common.Amount          `json:"refundedTotal"`
  GrandTotal            common.Amount          `json:"grandTotal"`
  ChannelName           string          `json:"channelName"`
  ExternalOrderReference string         `json:"externalOrderReference"`
  FulfilledOn           string          `json:"fulfilledOn"`
  PriceTaxInterpretation string         `json:"priceTaxInterpretation"`
}

type LineItem struct {
  ID             string            `json:"id"`
  VariantID      string            `json:"variantId,omitempty"`
  SKU            string            `json:"sku,omitempty"`
  Weight         float64           `json:"weight,omitempty"`
  Width          float64           `json:"width,omitempty"`
  Length         float64           `json:"length,omitempty"`
  Height         float64           `json:"height,omitempty"`
  ProductID      string            `json:"productId,omitempty"`
  ProductName    string            `json:"productName"`
  Quantity       int               `json:"quantity"`
  UnitPricePaid  common.Amount            `json:"unitPricePaid"`
  VariantOptions []VariantOption   `json:"variantOptions,omitempty"`
  Customizations []Customization   `json:"customizations,omitempty"`
  ImageURL       string            `json:"imageUrl,omitempty"`
  LineItemType   string            `json:"lineItemType"`
}

type VariantOption struct {
  Value      string `json:"value"`
  OptionName string `json:"optionName"`
}

type Customization struct {
  Label string `json:"label"`
  Value string `json:"value"`
}

type Note struct {
  Content string `json:"content"`
}

type ShippingLine struct {
  Method string `json:"method"`
  Amount common.Amount `json:"amount"`
}

type DiscountLine struct {
  Description string `json:"description,omitempty"`
  Name        string `json:"name"`
  Amount      common.Amount `json:"amount"`
  PromoCode   string `json:"promoCode,omitempty"`
}

type FormSubmission struct {
  Label string `json:"label"`
  Value string `json:"value"`
}

type Fulfillment struct {
  ShipDate      string `json:"shipDate"`
  CarrierName   string `json:"carrierName"`
  Service       string `json:"service"`
  TrackingNumber string `json:"trackingNumber"`
  TrackingURL   string `json:"trackingUrl"`
}

type Address struct {
  FirstName   string `json:"firstName"`
  LastName    string `json:"lastName"`
  Address1    string `json:"address1"`
  Address2    string `json:"address2,omitempty"`
  City        string `json:"city"`
  State       string `json:"state"`
  PostalCode  string `json:"postalCode"`
  CountryCode string `json:"countryCode"`
  Phone       string `json:"phone"`
}
