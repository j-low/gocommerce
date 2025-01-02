package products

import "time"

const (
	ProductsAPIVersion = "1.0"
)

type GetStorePagesRequest struct {
	Cursor string `json:"cursor"`
}

type GetStorePagesResponse struct {
	StorePages []StorePage `json:"storePages"`
	Pagination Pagination  `json:"pagination"`
}

type GetProductsResponse struct {
	Products   []Product   `json:"products"`
	Pagination Pagination  `json:"pagination"`
}

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

type UpdateVariantRequest struct {
	ProductID          string                  `json:"-"`
	VariantID          string                  `json:"-"`
	SKU                string                  `json:"sku,omitempty"`
	Pricing            *Pricing                `json:"pricing,omitempty"`
	Attributes         map[string]string       `json:"attributes,omitempty"`
	ShippingMeasurements *ShippingMeasurements `json:"shippingMeasurements,omitempty"`
}

type UpdateVariantResponse struct {
	ID                  string                  `json:"id"`
	SKU                 string                  `json:"sku"`
	Pricing             Pricing                 `json:"pricing"`
	Stock               *Stock                  `json:"stock,omitempty"`
	Attributes          map[string]string       `json:"attributes"`
	ShippingMeasurements *ShippingMeasurements  `json:"shippingMeasurements"`
	Image               *ProductImage           `json:"image,omitempty"`
}

type DownloadProductImageRequest struct {
	ProductID string
	PartID    string
	URL       string 
}

type DownloadProductImageResponse struct {
	S3Key string `json:"s3Key"`
}

type UploadProductImageRequest struct {
	ProductID string
	S3Key     string
}

type UploadProductImageResponse struct {
	ImageID string `json:"imageId"`
}

type VariantImageUpdateData struct {
	Product Product
	SKU string
	ImageID   string
	S3Key     string
}

type SetVariantImageRequest struct {
	ProductID string `json:"-"`
	VariantID string `json:"-"`
	ImageID   string `json:"imageId"`
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

type Pagination struct {
	HasNextPage    bool   `json:"hasNextPage"`
	NextPageCursor string `json:"nextPageCursor"`
	NextPageURL    string `json:"nextPageUrl"`
}

type SEOOptions struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Pricing struct {
	BasePrice Amount `json:"basePrice"`  
	OnSale    bool   `json:"onSale,omitempty"`
	SalePrice *Amount `json:"salePrice,omitempty"`
}

type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
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
