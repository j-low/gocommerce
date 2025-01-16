package profiles

import "github.com/NuvoCodeTechnologies/gocommerce/common"

const (
	ProfilesAPIVersion = "1.0"
)

type RetrieveAllProfilesResponse struct {
  Profiles   []Profile   `json:"profiles"`
  Pagination common.Pagination  `json:"pagination"`
}

type RetrieveSpecificProfilesResponse struct {
  Profiles []Profile `json:"profiles"`
}

type Profile struct {
	ID                  string               `json:"id"`
	FirstName           string               `json:"firstName"`
	LastName            string               `json:"lastName"`
	Email               string               `json:"email"`
	HasAccount          bool                 `json:"hasAccount"`
	IsCustomer          bool                 `json:"isCustomer"`
	CreatedOn           string               `json:"createdOn"`
	Address             *common.Address      `json:"address,omitempty"`
	AcceptsMarketing    bool                 `json:"acceptsMarketing"`
	TransactionsSummary *TransactionsSummary `json:"transactionsSummary,omitempty"`
}

type TransactionsSummary struct {
	FirstOrderSubmittedOn    *string        `json:"firstOrderSubmittedOn,omitempty"`
	LastOrderSubmittedOn     *string        `json:"lastOrderSubmittedOn,omitempty"`
	OrderCount               int            `json:"orderCount"`
	TotalOrderAmount         *common.Amount `json:"totalOrderAmount,omitempty"`
	TotalRefundAmount        *common.Amount `json:"totalRefundAmount,omitempty"`
	FirstDonationSubmittedOn *string        `json:"firstDonationSubmittedOn,omitempty"`
	LastDonationSubmittedOn  *string        `json:"lastDonationSubmittedOn,omitempty"`
	DonationCount            int            `json:"donationCount"`
	TotalDonationAmount      *common.Amount `json:"totalDonationAmount,omitempty"`
}
