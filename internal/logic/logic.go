package logic

import (
	"encoding/csv"
	"fmt"
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
	var apidata [][]string
	for idx := 1; idx < len(fields); idx++ {
		parts := strings.Split(fields[idx][5], "/")
		id := parts[len(parts)-1]
		data, err := l.service.FindInPep1(ctx, id, fields[idx][2], fields[idx][3], fields[idx][16])
		if err != nil {
			return nil, err
		}
		if data.Emails == nil {

			suspect := []string{fields[idx][2], fields[idx][3], fields[idx][16]}
			_ = append(apidata, suspect)
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

	if ctx.PostForm("responseType") == "json" {
		filename, err = utils.PayloadToJSON(resp, "data/req.csv", ctx.PostForm("email"),"scan")
		if err != nil {
			return nil, err
		}
	}

	if ctx.PostForm("responseType") == "csv" {
		filename, err = utils.PayloadToCSV(resp, "data/req.csv", ctx.PostForm("email"), "scan")
		if err != nil {
			return nil, err
		}
	}

	utils.SendFileToWebhook(os.Getenv("WEBHOOK_URL"), filename, ctx.PostForm("email"), ctx.PostForm("discordUsername"),ctx.PostForm("responseFormat"))
	return resp, nil
}

func (l Logic) ChangeWebhook(url string) error {
	err := os.Setenv("WEBHOOK_URL", url)
	if err != nil {
		fmt.Println("Error setting WEBHOOK_URL")
		return err
	}

	return nil
}

func (l Logic) GetPersonalEmail(file *multipart.FileHeader, ctx *gin.Context) ([]models.Payload, error) {
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
	var apidata [][]string

	for idx := 1; idx < len(fields); idx++ {
		parts := strings.Split(fields[idx][5], "/")
		id := parts[len(parts)-1]

		data, err := l.service.FindInPep1(ctx, id, fields[idx][2], fields[idx][3], fields[idx][16])
		if err != nil {
			return nil, err
		}

		if data.Emails == nil {

			suspect := []string{fields[idx][2], fields[idx][3], fields[idx][16]}
			_ = append(apidata, suspect)
		}

		resp = append(resp, models.Payload{
			Emails:    data.Emails,
			Telephone: data.Telephone,
		})

	}

	if ctx.PostForm("responseType") == "json" {
		filename, err = utils.PayloadToJSON(resp, "data/req.csv", ctx.PostForm("email"),"personal")
		if err != nil {
			return nil, err
		}
	}

	if ctx.PostForm("responseType") == "csv" {
		filename, err = utils.PayloadToCSV(resp, "data/req.csv", ctx.PostForm("email"), "personal")
		if err != nil {
			return nil, err
		}
	}

	utils.SendFileToWebhook(os.Getenv("WEBHOOK_URL"), filename, ctx.PostForm("email"), ctx.PostForm("discordUsername"),ctx.PostForm("responseFormat"))
	return resp, nil	
}

func (l Logic) GetProfessionalEmail(file *multipart.FileHeader, ctx *gin.Context) ([]models.Payload, error) {
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
	var apidata [][]string

	for idx := 1; idx < len(fields); idx++ {
		parts := strings.Split(fields[idx][5], "/")
		id := parts[len(parts)-1]

		data, err := l.service.GetProfessionalEmails(ctx, id, fields[idx][2], fields[idx][3], fields[idx][16])
		if err != nil {
			return nil, err
		}

		if data.Emails == nil {

			suspect := []string{fields[idx][2], fields[idx][3], fields[idx][16]}
			_ = append(apidata, suspect)
		}

		resp = append(resp, models.Payload{
			Emails:    data.Emails,
			Telephone: data.Telephone,
		})

	}

	if ctx.PostForm("responseType") == "json" {
		filename, err = utils.PayloadToJSON(resp, "data/req.csv", ctx.PostForm("email"),"professional")
		if err != nil {
			return nil, err
		}
	}

	if ctx.PostForm("responseType") == "csv" {
		filename, err = utils.PayloadToCSV(resp, "data/req.csv", ctx.PostForm("email"), "professional")
		if err != nil {
			return nil, err
		}
	}

	utils.SendFileToWebhook(os.Getenv("WEBHOOK_URL"), filename, ctx.PostForm("email"), ctx.PostForm("discordUsername"), ctx.PostForm("responseFormat"))
	return resp, nil	
}


func (l Logic) GetBothEmails(file *multipart.FileHeader, ctx *gin.Context) ([]models.Payload, error) {
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
	var apidata [][]string
	for idx := 1; idx < len(fields); idx++ {
		parts := strings.Split(fields[idx][5], "/")
		id := parts[len(parts)-1]
		data, err := l.service.FindInPep1(ctx, id, fields[idx][2], fields[idx][3], fields[idx][16])
		if err != nil {
			return nil, err
		}
		if data.Emails == nil {

			suspect := []string{fields[idx][2], fields[idx][3], fields[idx][16]}
			_ = append(apidata, suspect)
		}
		//
		resp = append(resp, models.Payload{
			Emails:    data.Emails,
			Telephone: data.Telephone,
		})
		//2 3 5 16
	}

	if ctx.PostForm("responseType") == "json" {
		filename, err = utils.PayloadToJSONforFiltering(resp, "data/req.csv", ctx.PostForm("email"))
		if err != nil {
			return nil, err
		}
	}

	if ctx.PostForm("responseType") == "csv" {
		filename, err = utils.PayloadToCSVforFiltering(resp, "data/req.csv", ctx.PostForm("email"))
		if err != nil {
			return nil, err
		}
	}

	utils.SendFileToWebhook(os.Getenv("WEBHOOK_URL"), filename, ctx.PostForm("email"), ctx.PostForm("discordUsername"),ctx.PostForm("responseFormat"))
	return resp, nil
}

func (l Logic) GetByLIID(ctx *gin.Context,liid string) (models.Payload, error) {
	data, err := l.service.FindInPep1(ctx, liid, "", "", "")
	if err != nil {
		return models.Payload{}, err
	}

	return models.Payload{
		Emails:    data.Emails,
		Telephone: data.Telephone,
	}, nil
}

func (l Logic) GetMultipleByLIID(ctx *gin.Context,liids []string) ([]models.Payload, error) {
	var resp []models.Payload
	for _,liid := range liids {
		data, err := l.service.FindInPep1(ctx, liid, "", "", "")
		if err != nil {
			return nil, err
		}

		resp = append(resp, models.Payload{
			Emails:    data.Emails,
			Telephone: data.Telephone,
		})
	}

	return resp, nil
}

func (l Logic) GetPersonalEmailByliid(ctx *gin.Context,liid string) (models.Payload, error) {
	data, err := l.service.GetPersonalEmail(ctx, liid, "", "", "")
	if err != nil {
		return models.Payload{}, err
	}

	return models.Payload{
		Emails:    data.Emails,
		Telephone: data.Telephone,
	}, nil
}

func (l Logic) GetProfessionalEmailsByliid(ctx *gin.Context,liid string) (models.Payload, error) {
	data, err := l.service.GetProfessionalEmails(ctx, liid, "", "", "")
	if err != nil {
		return models.Payload{}, err
	}

	return models.Payload{
		Emails:    data.Emails,
		Telephone: data.Telephone,
	}, nil
}