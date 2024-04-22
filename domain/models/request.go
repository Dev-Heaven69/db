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
}

type ApiResponse struct {
	State bool `json:"state"`
	Email string `json:"email"`
}

type Response struct {
	Data []Payload `json:"data"`
	ResquesteeEmail string `json:"requesteeEmail"`
}