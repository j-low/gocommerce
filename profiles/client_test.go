package profiles

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/j-low/gocommerce/common"
)

func TestRetrieveAllProfiles(t *testing.T) {
	tests := []struct {
		name        string
		params      common.QueryParams
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful retrieval",
			params:     common.QueryParams{},
			mockStatus: http.StatusOK,
			mockResp: `{
                "profiles": [
                    {
                        "id": "profile-123",
                        "firstName": "John",
                        "lastName": "Doe",
                        "email": "john@example.com",
                        "phone": "123-456-7890",
                        "addresses": [
                            {
                                "firstName": "John",
                                "lastName": "Doe",
                                "address1": "123 Main St",
                                "city": "New York",
                                "state": "NY",
                                "countryCode": "US",
                                "postalCode": "10001"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "nextCursor": "next-page"
                }
            }`,
			wantErr: false,
		},
		{
			name: "with cursor",
			params: common.QueryParams{
				Cursor: "next-page",
			},
			mockStatus: http.StatusOK,
			mockResp: `{
                "profiles": [],
                "pagination": {}
            }`,
			wantErr: false,
		},
		{
			name:        "server error",
			params:      common.QueryParams{},
			mockStatus:  http.StatusInternalServerError,
			mockResp:    `{"type":"ERROR","message":"Internal server error"}`,
			wantErr:     true,
			errContains: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %s", auth)
				}

				if tt.params.Cursor != "" {
					cursor := r.URL.Query().Get("cursor")
					if cursor != tt.params.Cursor {
						t.Errorf("expected cursor %s, got %s", tt.params.Cursor, cursor)
					}
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			config := &common.Config{
				APIKey:    "test-key",
				Client:    server.Client(),
				UserAgent: "test-agent",
				BaseURL:   server.URL,
			}

			resp, err := RetrieveAllProfiles(context.Background(), config, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAllProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if !tt.wantErr && resp == nil {
				t.Error("expected non-nil response when no error")
			}
		})
	}
}

func TestRetrieveSpecificProfiles(t *testing.T) {
	tests := []struct {
		name        string
		profileIDs  []string
		mockStatus  int
		mockResp    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful retrieval",
			profileIDs: []string{"5ede64464561317766bdc632", "5fc7b8f18b4dc33dc818aeba"},
			mockStatus: http.StatusOK,
			mockResp: `{
                "profiles": [
                    {
                        "id": "5ede64464561317766bdc632",
                        "firstName": "Gregory",
                        "lastName": "Jones",
                        "email": "gregory_jones@example.com",
                        "hasAccount": false,
                        "isCustomer": true,
                        "createdOn": "2020-06-08T16:16:06.573518Z",
                        "acceptsMarketing": false,
                        "address": {
                            "firstName": "Gregory",
                            "lastName": "Jones",
                            "address1": "450 North End Avenue",
                            "city": "New York",
                            "state": "NY",
                            "countryCode": "US",
                            "postalCode": "10282",
                            "phone": "5553334444"
                        }
                    }
                ]
            }`,
			wantErr: false,
		},
		{
			name:        "too many profiles",
			profileIDs:  make([]string, 51), // More than 50 IDs
			mockStatus:  http.StatusBadRequest,
			mockResp:    `{"type":"INVALID_REQUEST_ERROR","message":"Cannot request more than 50 profiles"}`,
			wantErr:     true,
			errContains: "Cannot request more than 50 profiles",
		},
		{
			name:        "profile not found",
			profileIDs:  []string{"non-existent"},
			mockStatus:  http.StatusNotFound,
			mockResp:    `{"type":"INVALID_REQUEST_ERROR","message":"Profile not found","subtype":"INVALID_ARGUMENT"}`,
			wantErr:     true,
			errContains: "Profile not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/1.0/profiles/%s", strings.Join(tt.profileIDs, ","))
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %s", auth)
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			config := &common.Config{
				APIKey:    "test-key",
				Client:    server.Client(),
				UserAgent: "test-agent",
				BaseURL:   server.URL,
			}

			resp, err := RetrieveSpecificProfiles(context.Background(), config, tt.profileIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpecificProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if !tt.wantErr {
				if resp == nil || len(resp.Profiles) == 0 {
					t.Error("expected non-nil response with profiles when no error")
				}
			}
		})
	}
}
