package common

import "net/http"

type Config struct {
	APIKey string
	UserAgent string
	Client *http.Client
}

type APIError struct {
	Type string
	Subtype string
	Message string
	Detail string
}