package common

import (
	"net/http"

	"github.com/google/uuid"
)

type Config struct {
	APIKey string
	UserAgent string
	Client *http.Client
	IdempotencyKey  *uuid.UUID
}

type QueryParams struct {
	Cursor        string
	Filter        string
	ModifiedAfter string
	ModifiedBefore string
	SortDirection string
	SortField     string
	Status        string
}

type Pagination struct {
	HasNextPage    bool   `json:"hasNextPage"`
	NextPageCursor string `json:"nextPageCursor"`
	NextPageURL    string `json:"nextPageUrl"`
}

type APIError struct {
	Type string
	Subtype string
	Message string
	Detail string
}

type Address struct {
  FirstName  string `json:"firstName"`
  LastName   string `json:"lastName"`
  Address1   string `json:"address1"`
  Address2   string `json:"address2,omitempty"`
  City       string `json:"city"`
  State      string `json:"state"`
  PostalCode string `json:"postalCode"`
  CountryCode string `json:"countryCode"`
  Phone      string `json:"phone"`
}

type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}