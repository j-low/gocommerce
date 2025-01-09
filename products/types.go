package products

import (
	"time"

	"github.com/NuvoCodeTechnologies/gocommerce/common"
)

const (
	ProductsAPIVersion = "1.0"
)

type CreateProductRequest struct {
	Type              string            `json:"type"`
	StorePageID       string            `json:"storePageId"`
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
	URLSlug           string            `json:"urlSlug,omitempty"`
	Tags              []string          `json:"tags,omitempty"`
	IsVisible         bool              `json:"isVisible"`
	VariantAttributes []string          `json:"variantAttributes,omitempty"`
	Variants          []ProductVariant  `json:"variants"`
}

type CreateProductResponse struct {
	ID                string           `json:"id"`
	Type              string           `json:"type"`
	StorePageID       string           `json:"storePageId"`
	Name              string           `json:"name,omitempty"`
	Description       string           `json:"description,omitempty"`
	URL               string           `json:"url"`
	URLSlug           string           `json:"urlSlug"`
	Tags              []string         `json:"tags,omitempty"`
	IsVisible         bool             `json:"isVisible"`
	VariantAttributes []string         `json:"variantAttributes"`
	Variants          []ProductVariant `json:"variants"`
	Images            []string         `json:"images,omitempty"`
	CreatedOn         string           `json:"createdOn"`
}

type CreateProductVariantRequest struct {
  ProductID          string                  `json:"-"`
  SKU                string                  `json:"sku"`
  Pricing            Pricing                 `json:"pricing"`
  Stock              *Stock                  `json:"stock,omitempty"`
  Attributes         map[string]string       `json:"attributes,omitempty"`
  ShippingMeasurements *ShippingMeasurements `json:"shippingMeasurements,omitempty"`
}

type CreateProductVariantResponse struct {
  ID                 string                  `json:"id"`
  SKU                string                  `json:"sku"`
  Pricing            Pricing                 `json:"pricing"`
  Stock              *Stock                  `json:"stock,omitempty"`
  Attributes         map[string]string       `json:"attributes,omitempty"`
  ShippingMeasurements *ShippingMeasurements `json:"shippingMeasurements,omitempty"`
  Image              *ProductImage           `json:"image,omitempty"`
}

type UploadProductImageRequest struct {
  ProductID string
  FilePath  string
}

type UploadProductImageResponse struct {
  ImageID string `json:"id"`
}

type RetrieveAllStorePagesResponse struct {
	StorePages []StorePage `json:"storePages"`
	Pagination common.Pagination  `json:"pagination"`
}

type RetrieveAllProductsRequest struct {
	ModifiedAfter 	string `json:"modifiedAfter"` 
	ModifiedBefore 	string `json:"modifiedBefore"`
	Type          	string `json:"type"`
}

type RetrieveAllProductsResponse struct {
	Products   []Product   `json:"products"`
	Pagination common.Pagination  `json:"pagination"`
}

type RetrieveSpecificProductsRequest struct {
  ProductIDs []string `json:"productIds"` // List of specific product IDs to retrieve
}

type RetrieveSpecificProductsResponse struct {
  Products []Product `json:"products"`
}

type GetProductImageUploadStatusResponse struct {
  Status string `json:"status"`
}

type AssignProductImageToVariantRequest struct {
	ProductID string `json:"-"`
	VariantID string `json:"-"`
	ImageID   string `json:"imageId"`
}

type ReorderProductImageRequest struct {
  ProductID   string  `json:"-"`
  ImageID     string  `json:"-"`
  AfterImageID *string `json:"afterImageId,omitempty"`
}

type UpdateProductRequest struct {
	Name              *string           `json:"name,omitempty"`
	Description       *string           `json:"description,omitempty"`
	URLSlug           *string           `json:"urlSlug,omitempty"`
	Tags              []string          `json:"tags,omitempty"`
	IsVisible         *bool             `json:"isVisible,omitempty"`
	VariantAttributes []string          `json:"variantAttributes,omitempty"`
	SEOOptions        *SEOOptions       `json:"seoOptions,omitempty"`
}

type UpdateProductResponse struct {
	ID                string           `json:"id"`
	Type              string           `json:"type"`
	StorePageID       string           `json:"storePageId"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	URL               string           `json:"url"`
	URLSlug           string           `json:"urlSlug"`
	Tags              []string         `json:"tags"`
	IsVisible         bool             `json:"isVisible"`
	SEOOptions        SEOOptions       `json:"seoOptions"`
	VariantAttributes []string         `json:"variantAttributes"`
	Variants          []ProductVariant `json:"variants"`
	Images            []ProductImage   `json:"images"`
	CreatedOn         string           `json:"createdOn"`
	ModifiedOn        string           `json:"modifiedOn"`
}

type UpdateProductVariantRequest struct {
	ProductID          string                  `json:"-"`
	VariantID          string                  `json:"-"`
	SKU                string                  `json:"sku,omitempty"`
	Pricing            *Pricing                `json:"pricing,omitempty"`
	Attributes         map[string]string       `json:"attributes,omitempty"`
	ShippingMeasurements *ShippingMeasurements `json:"shippingMeasurements,omitempty"`
}

type UpdateProductVariantResponse struct {
	ID                  string                  `json:"id"`
	SKU                 string                  `json:"sku"`
	Pricing             Pricing                 `json:"pricing"`
	Stock               *Stock                  `json:"stock,omitempty"`
	Attributes          map[string]string       `json:"attributes"`
	ShippingMeasurements *ShippingMeasurements  `json:"shippingMeasurements"`
	Image               *ProductImage           `json:"image,omitempty"`
}

type UpdateProductImageRequest struct {
  ProductID string `json:"-"`
  ImageID   string `json:"-"`
  AltText   string `json:"altText"`
}

type UpdateProductImageResponse struct {
  ID               string    `json:"id"`
  AltText          string    `json:"altText"`
  URL              string    `json:"url"`
  OriginalSize     ImageSize `json:"originalSize"`
  AvailableFormats []string  `json:"availableFormats"`
}

type StorePage struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	IsEnabled bool   `json:"isEnabled"`
}

type Product struct {
	ID                string         `json:"id"`
	Type              string         `json:"type"`
	StorePageID       string         `json:"storePageId"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	URL               string         `json:"url"`
	URLSlug           string         `json:"urlSlug"`
	Tags              []string       `json:"tags"`
	IsVisible         bool           `json:"isVisible"`
	SEOOptions        SEOOptions     `json:"seoOptions"`
	VariantAttributes []string       `json:"variantAttributes"`
	Variants          []ProductVariant `json:"variants"`
	Images            []ProductImage `json:"images"`
	Pricing           *Pricing       `json:"pricing,omitempty"`
	DigitalGood       *DigitalGood   `json:"digitalGood,omitempty"`
	CreatedOn         time.Time      `json:"createdOn"`
	ModifiedOn        time.Time      `json:"modifiedOn"`
}

type ProductVariant struct {
	ID                 string                `json:"id"`
	SKU                string                `json:"sku"`
	Pricing            Pricing               `json:"pricing"`
	Stock              *Stock                `json:"stock,omitempty"`
	Attributes         map[string]string     `json:"attributes,omitempty"`
	ShippingMeasurements *ShippingMeasurements `json:"shippingMeasurements,omitempty"`
}

type SEOOptions struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Pricing struct {
	BasePrice common.Amount `json:"basePrice"`  
	OnSale    bool   `json:"onSale,omitempty"`
	SalePrice *common.Amount `json:"salePrice,omitempty"`
}

type DigitalGood struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type Stock struct {
	Quantity  int  `json:"quantity,omitempty"`
	Unlimited bool `json:"unlimited,omitempty"`
}

type ShippingMeasurements struct {
	Weight     *Weight     `json:"weight,omitempty"`
	Dimensions *Dimensions `json:"dimensions,omitempty"`
}

type Weight struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

type Dimensions struct {
	Unit   string  `json:"unit"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type ProductImage struct {
	ID               string       `json:"id"`
	AltText          string       `json:"altText"`
	URL              string       `json:"url"`
	OriginalSize     ImageSize    `json:"originalSize"`
	AvailableFormats []string     `json:"availableFormats"`
}

type ImageSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
