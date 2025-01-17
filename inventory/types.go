package inventory

import "github.com/j-low/gocommerce/common"

const (
	InventoryAPIVersion = "1.0"
)

type RetrieveAllInventoryResponse struct {
	Inventory  []InventoryRecord `json:"inventory"`
	Pagination common.Pagination `json:"pagination"`
}

type RetrieveSpecificInventoryResponse struct {
	Inventory []InventoryRecord `json:"inventory"`
}

type AdjustStockQuantitiesRequest struct {
	IncrementOperations    []QuantityOperation `json:"incrementOperations,omitempty"`
	DecrementOperations    []QuantityOperation `json:"decrementOperations,omitempty"`
	SetFiniteOperations    []QuantityOperation `json:"setFiniteOperations,omitempty"`
	SetUnlimitedOperations []string            `json:"setUnlimitedOperations,omitempty"`
}

type InventoryRecord struct {
	VariantID   string `json:"variantId"`
	SKU         string `json:"sku"`
	Descriptor  string `json:"descriptor"`
	IsUnlimited bool   `json:"isUnlimited"`
	Quantity    int    `json:"quantity"`
}

type QuantityOperation struct {
	VariantID string `json:"variantId"`
	Quantity  int    `json:"quantity"`
}
