package models

import "mime/multipart"

type CSVRequest struct {
	CsvFile *multipart.FileHeader `form:"csv" binding:"required"`
}

type ApiResponse struct {
	State bool `json:"state"`
	Email string `json:"email"`
}