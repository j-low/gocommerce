package inventory

import "github.com/NuvoCodeTechnologies/gocommerce/common"

const (
	InventoryAPIVersion = "1.0"
)

type RetrieveAllInventoryResponse struct {
  Inventory []InventoryRecord `json:"inventory"`
  Pagination common.Pagination `json:"pagination"`
}

type RetrieveSpecificInventoryRequest struct {
  InventoryIDs []string `json:"inventoryIds"`
}

type RetrieveSpecificInventoryResponse struct {
  Inventory []InventoryRecord `json:"inventory"`
}

type AdjustStockQuantitiesRequest struct {
  Adjustments []StockAdjustment `json:"adjustments"`
}

type AdjustStockQuantitiesResponse struct {
  Results []AdjustmentResult `json:"results"`
}

type InventoryRecord struct {
  ID        string `json:"id"`
  ProductID string `json:"productId"`
  VariantID string `json:"variantId"`
  Stock     int    `json:"stock"`
  Unlimited bool   `json:"unlimited"`
}

type StockAdjustment struct {
  InventoryID string `json:"inventoryId"`
  Delta       int    `json:"delta"`
}

type AdjustmentResult struct {
  InventoryID string `json:"inventoryId"`
  Success     bool   `json:"success"`
  Error       string `json:"error,omitempty"`
}
