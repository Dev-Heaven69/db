package models

import (
	"mime/multipart"
	"net/mail"
)

type Email string
type ResponseType string

func (rt *ResponseType) IsValid() bool {
	return *rt == "json" || *rt == "csv"
}

func (e *Email) IsValid() bool {
	_,err := mail.ParseAddress(string(*e))
	return err == nil
}

type Request struct {
	CsvFile *multipart.FileHeader `form:"csv" binding:"required"`
	ResponseType ResponseType `form:"responseType" binding:"required"`
	DiscordUsername string `form:"discordUsername" binding:"required"`
	Email Email `form:"email" binding:"required"`
	ResponseFormat string `form:"responseFormat" binding:"required"`
}

type ScanRequest struct {
	CsvFile *multipart.FileHeader `form:"csv" binding:"required"`
	DiscordUsername string `form:"discordUsername" binding:"required"`
	Email Email `form:"email" binding:"required"`
	WantedFields string `form:"wantedFields" binding:"required"`
}

type ApiResponse struct {
	State bool `json:"state"`
	Email string `json:"email"`
}

type ChangeWebhookRequest struct {
	URL string `json:"url" binding:"required"`
}

type GetOneByLIIDRequest struct {
	Liid string `json:"liid" binding:"required"`
}	


type GetMultipleByLIIDRequest struct {
	Liids []string `json:"liids" binding:"required"`
}
