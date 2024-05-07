package models

type Pep1Response struct {
	Emails    []string `bson:"e,omitempty" `
	Telephone []string `bson:"t,omitempty"`
}

type Payload struct {
	Emails    []string `json:"emails"`
	Telephone []string `json:"phoneNumbers"`
}
