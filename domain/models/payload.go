package models

type Pep1Response struct {
	Emails    []string `bson:"e,omitempty" bson:"emails,omitempty"`
	Telephone []string `bson:"t,omitempty" bson:"phone_numbers,omitempty"`
}

// type UseResponse struct {
// 	Emails    []string `json:"emails" bson:"emails,omitempty"`
// 	Telephone []string `json:"phone_numbers" bson:"phone_numbers,omitempty"`
// }

type Payload struct {
	Emails    []string `json:"emails"`
	Telephone []string `json:"phoneNumbers"`
}
