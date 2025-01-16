package transactions

import "github.com/j-low/gocommerce/common"

const (
	TransactionsAPIVersion = "1.0"
)

type RetrieveAllTransactionsResponse struct {
	Documents   []Document         `json:"documents"`
	Pagination  common.Pagination  `json:"pagination"`
}

type RetrieveSpecificTransactionsResponse struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	ID                   string             `json:"id"`
	CreatedOn            string             `json:"createdOn"`
	ModifiedOn           string             `json:"modifiedOn"`
	CustomerEmail        *string            `json:"customerEmail,omitempty"`
	SalesOrderID         *string            `json:"salesOrderId,omitempty"`
	Voided               bool               `json:"voided"`
	TotalSales           common.Amount      `json:"totalSales"`
	TotalNetSales        common.Amount      `json:"totalNetSales"`
	TotalNetShipping     common.Amount      `json:"totalNetShipping"`
	TotalTaxes           common.Amount      `json:"totalTaxes"`
	Total                common.Amount      `json:"total"`
	TotalNetPayment      common.Amount      `json:"totalNetPayment"`
	Payments             []Payment          `json:"payments"`
	SalesLineItems       []SalesLineItem    `json:"salesLineItems"`
	Discounts            []Discount         `json:"discounts"`
	ShippingLineItems    []ShippingLineItem `json:"shippingLineItems"`
	PaymentGatewayError  *string            `json:"paymentGatewayError,omitempty"`
}

type Payment struct {
	ID                         string            `json:"id"`
	Amount                     common.Amount     `json:"amount"`
	RefundedAmount             common.Amount     `json:"refundedAmount"`
	NetAmount                  common.Amount     `json:"netAmount"`
	CreditCardType             *string           `json:"creditCardType,omitempty"`
	Provider                   string            `json:"provider"`
	Refunds                    []Refund          `json:"refunds"`
	ProcessingFees             []ProcessingFee    `json:"processingFees"`
	GiftCardID                 *string           `json:"giftCardId,omitempty"`
	PaidOn                     string            `json:"paidOn"`
	ExternalTransactionID      string            `json:"externalTransactionId"`
	ExternalTransactionProperties []interface{}  `json:"externalTransactionProperties"`
	ExternalCustomerID         *string           `json:"externalCustomerId,omitempty"`
}

type Refund struct {
	ID                    string        `json:"id"`
	Amount                common.Amount `json:"amount"`
	RefundedOn            string        `json:"refundedOn"`
	ExternalTransactionID string        `json:"externalTransactionId"`
}

type ProcessingFee struct {
	ID                         string            `json:"id"`
	Amount                     common.Amount     `json:"amount"`
	AmountGatewayCurrency      common.Amount     `json:"amountGatewayCurrency"`
	ExchangeRate               string            `json:"exchangeRate"`
	RefundedAmount             common.Amount     `json:"refundedAmount"`
	RefundedAmountGatewayCurrency common.Amount `json:"refundedAmountGatewayCurrency"`
	NetAmount                  common.Amount     `json:"netAmount"`
	NetAmountGatewayCurrency   common.Amount     `json:"netAmountGatewayCurrency"`
	FeeRefunds                 []FeeRefund       `json:"feeRefunds"`
}

type FeeRefund struct {
	ID                         string        `json:"id"`
	Amount                     common.Amount `json:"amount"`
	AmountGatewayCurrency      common.Amount `json:"amountGatewayCurrency"`
	ExchangeRate               string        `json:"exchangeRate"`
	RefundedOn                 string        `json:"refundedOn"`
	ExternalTransactionID      string        `json:"externalTransactionId"`
}

type SalesLineItem struct {
	ID             string        `json:"id"`
	DiscountAmount common.Amount `json:"discountAmount"`
	TotalSales     common.Amount `json:"totalSales"`
	TotalNetSales  common.Amount `json:"totalNetSales"`
	Total          common.Amount `json:"total"`
	Taxes          []Tax         `json:"taxes"`
}

type Tax struct {
	Amount       common.Amount `json:"amount"`
	Rate         string        `json:"rate"`
	Name         string        `json:"name"`
	Jurisdiction string        `json:"jurisdiction,omitempty"`
	Description  string        `json:"description,omitempty"` 
}

type Discount struct {
	Description string        `json:"description,omitempty"`
	Name        string        `json:"name"`
	Amount      common.Amount `json:"amount"`
}

type ShippingLineItem struct {
	ID             string        `json:"id"`
	Amount         common.Amount `json:"amount"`
	DiscountAmount common.Amount `json:"discountAmount"`
	NetAmount      common.Amount `json:"netAmount"`
	Description     string        `json:"description"`
	Taxes           []Tax         `json:"taxes"`
}
