package mocks

import "net/http"

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}
