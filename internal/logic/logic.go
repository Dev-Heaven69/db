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

func (l Logic) FindPep1(file *multipart.FileHeader, ctx *gin.Context) (models.Response, error) {
	uploadPath := "./data/"
	filename := "req.csv"

	filepath := uploadPath + filename

	err := ctx.SaveUploadedFile(file, filepath)
	if err != nil {
		return models.Response{}, err
	}
	csvFile, err := os.Open(filepath)
	if err != nil {
		return models.Response{}, err
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
	var apidata [][]string
	for idx := 1; idx < len(fields); idx++ {
		parts := strings.Split(fields[idx][5], "/")
		id := parts[len(parts)-1]
		data, err := l.service.FindInPep1(ctx, id, fields[idx][2], fields[idx][3], fields[idx][16])
		if err != nil {
			return models.Response{}, err
		}
		if data.Emails == nil {

			suspect := []string{fields[idx][2], fields[idx][3], fields[idx][16]}
			apidata = append(apidata, suspect)
		}
		//
		resp = append(resp, models.Payload{
			Emails:    data.Emails,
			Telephone: data.Telephone,
		})
		//2 3 5 16
	}

	// fmt.Println("this is apidata",apidata)

	// limiter := rate.NewLimiter(5, 1)

	// datafromAPI := utils.QueryBulkRecords()
	// fmt.Println("this is datafromAPI", datafromAPI)
	var apiResponse models.Response
	apiResponse.Data = resp
	apiResponse.ResquesteeEmail = string(ctx.PostForm("email"))
	if ctx.PostForm("responseType") == "json"{
		filename,err = utils.WriteResponseToJson(apiResponse)
		if err != nil {
			return models.Response{}, err
		}
		utils.SendToWebhook("http://n8n.leadzenai.co/webhook/ewH5SNa0IhYTsyZi/webhook1/receive-json",filename, ctx.PostForm("responseType"))
	}
	if ctx.PostForm("responseType") == "csv"{
		filename,err = utils.PayloadToCSV(apiResponse, "data/req.csv")
		if err != nil {
			return models.Response{}, err
		}
		utils.SendToWebhook("http://n8n.leadzenai.co/webhook/ewH5SNa0IhYTsyZi/webhook1/receive-csv",filename, ctx.PostForm("responseType"))	
	}
	return apiResponse, nil
}
