package models

import "mime/multipart"

type CSVRequest struct {
	CsvFile *multipart.FileHeader `form:"csv" binding:"required"`
}