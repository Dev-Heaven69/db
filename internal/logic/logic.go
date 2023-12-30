package logic

import (
	"encoding/csv"
	"mime/multipart"
	"os"
	"strings"

	"github.com/DevHeaven/db/domain/models"
	"github.com/DevHeaven/db/internal/dbi"
	"github.com/DevHeaven/db/internal/utils"
	"github.com/gin-gonic/gin"
)

type Logic struct {
	service dbi.Service
}

func ProvideLogic(service dbi.Service) Logic {
	return Logic{service: service}
}

func (l Logic) FindPep1(file *multipart.FileHeader, ctx *gin.Context) ([]models.Payload, error) {
	uploadPath := "./data/"
	filename := "req.csv"

	filepath := uploadPath + filename

	err := ctx.SaveUploadedFile(file, filepath)
	if err != nil {
		return nil, err
	}
	csvFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	fields := make([][]string, 0)
	reader := csv.NewReader(csvFile)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		fields = append(fields, record)
	}

	var resp []models.Payload
	for idx := 1; idx < len(fields); idx++ {
		parts := strings.Split(fields[idx][5], "/")
		id := parts[len(parts)-1]
		data, err := l.service.FindInPep1(ctx, id)
		if err != nil {
			return nil, err
		}
		resp = append(resp, models.Payload{
			Emails:    data.Emails,
			Telephone: data.Telephone,
		})
		//2 3 5 16
	}
	// utils.CreateCSV(resp, "data/response.csv")
	utils.PayloadToCSV(resp, "data/req.csv")
	utils.SendCSVToWebhook("http://n8n.leadzenai.co/webhook-test/ewH5SNa0IhYTsyZi/webhook1/receive-csv")
	return resp, nil
}
