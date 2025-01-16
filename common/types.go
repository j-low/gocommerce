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

type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}