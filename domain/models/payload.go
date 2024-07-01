package models

type DbResponse struct {
	Emails    []string `bson:"e,omitempty" `
	Telephone []string `bson:"t,omitempty"`
	FirstName string   `bson:"First Name"`
	LastName  string   `bson:"Last Name"`
	OrganizationName string `bson:"Organization Name"`
	OrganizationDomain string `bson:"Organization Domain"`
	LinkedInUrl string `bson:"linkedin"`
}

type Payload struct {
	Emails    []string `json:"emails"`
	Telephone []string `json:"phoneNumbers"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	OrganizationName string `json:"OrganizationName"`
	OrganizationDomain string `json:"organizationDomain"`
	LinkedInUrl string `json:"linkedInUrl"`
}
