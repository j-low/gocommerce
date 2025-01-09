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
  ID              string   `json:"id"`
  FirstName       string   `json:"firstName"`
  LastName        string   `json:"lastName"`
  Email           string   `json:"email"`
  HasAccount      bool     `json:"hasAccount"`
  IsCustomer      bool     `json:"isCustomer"`
  CreatedOn       string   `json:"createdOn"`
  Address         *Address `json:"address,omitempty"`
  AcceptsMarketing bool    `json:"acceptsMarketing"`
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
