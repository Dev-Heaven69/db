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

type CSVFileData struct {
	FirstName          []string
	LastName           []string
	OrganizationDomain []string
	Emails             []string
	PhoneNumbers       []string
	Liid               []string
	LinkedInURL        []string
	PersonalEmails     []string
	ProfessionalEmails []string
}

func (l Logic) ScanDB(file *multipart.FileHeader, ctx *gin.Context) ([]models.Payload, error) {
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
		parts := strings.Split(fields[idx][4], "/")
		id := parts[len(parts)-1]
		data, err := l.service.ScanDB(ctx, id, "liid")
		if err != nil {
			return nil, err
		}

		resp = append(resp, models.Payload{
			Emails:    data.Emails,
			Telephone: data.Telephone,
		})
	}

	if ctx.PostForm("responseType") == "json" {
		filename, err = utils.PayloadToJSON(resp, "data/req.csv", ctx.PostForm("email"), "scan")
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

	utils.SendFileToWebhook(os.Getenv("WEBHOOK_URL"), filename, ctx.PostForm("email"), ctx.PostForm("discordUsername"), ctx.PostForm("responseFormat"))
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

		data, err := l.service.ScanDB(ctx, id, "liid")
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
		filename, err = utils.PayloadToJSON(resp, "data/req.csv", ctx.PostForm("email"), "personal")
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

	utils.SendFileToWebhook(os.Getenv("WEBHOOK_URL"), filename, ctx.PostForm("email"), ctx.PostForm("discordUsername"), ctx.PostForm("responseFormat"))
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

		data, err := l.service.GetProfessionalEmails(ctx, id)
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
		filename, err = utils.PayloadToJSON(resp, "data/req.csv", ctx.PostForm("email"), "professional")
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
		data, err := l.service.ScanDB(ctx, id, "liid")
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

	utils.SendFileToWebhook(os.Getenv("WEBHOOK_URL"), filename, ctx.PostForm("email"), ctx.PostForm("discordUsername"), ctx.PostForm("responseFormat"))
	return resp, nil
}

func (l Logic) GetByLIID(ctx *gin.Context, liid string) (models.Payload, error) {
	data, err := l.service.ScanDB(ctx, liid, "liid")
	if err != nil {
		return models.Payload{}, err
	}

	return models.Payload{
		Emails:    data.Emails,
		Telephone: data.Telephone,
	}, nil
}

func (l Logic) GetMultipleByLIID(ctx *gin.Context, liids []string) ([]models.Payload, error) {
	var resp []models.Payload
	for _, liid := range liids {
		data, err := l.service.ScanDB(ctx, liid, "liid")
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

func (l Logic) GetPersonalEmailByliid(ctx *gin.Context, liid string) (models.Payload, error) {
	data, err := l.service.GetPersonalEmail(ctx, liid)
	if err != nil {
		return models.Payload{}, err
	}

	return models.Payload{
		Emails:    data.Emails,
		Telephone: data.Telephone,
	}, nil
}

func (l Logic) GetProfessionalEmailsByliid(ctx *gin.Context, liid string) (models.Payload, error) {
	data, err := l.service.GetProfessionalEmails(ctx, liid)
	if err != nil {
		return models.Payload{}, err
	}

	return models.Payload{
		Emails:    data.Emails,
		Telephone: data.Telephone,
	}, nil
}

func (l Logic) Test(file *multipart.FileHeader, ctx *gin.Context) (CSVFileData, error) {
	uploadPath := "./data/"
	filename := "req.csv"
	var err error

	filepath := uploadPath + filename

	err = ctx.SaveUploadedFile(file, filepath)
	if err != nil {
		return CSVFileData{}, err
	}

	csvFile, err := os.Open(filepath)
	if err != nil {
		return CSVFileData{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	headers, err := reader.Read()
	if err != nil {
		return CSVFileData{}, err
	}

	fields, err := reader.ReadAll()
	if err != nil {
		return CSVFileData{}, err
	}

	var csvDataStruct CSVFileData
	for _, record := range fields {
		for i, value := range record {
			switch headers[i] {
			case "First Name":
				csvDataStruct.FirstName = append(csvDataStruct.FirstName, value)
			case "Last Name":
				csvDataStruct.LastName = append(csvDataStruct.LastName, value)
			case "Organization Domain":
				csvDataStruct.OrganizationDomain = append(csvDataStruct.OrganizationDomain, value)
			case "Phone Number":
				csvDataStruct.PhoneNumbers = append(csvDataStruct.PhoneNumbers, value)
			case "LinkedIn":
				id := strings.Split(value, "/")
				csvDataStruct.Liid = append(csvDataStruct.Liid, id[len(id)-1])
				csvDataStruct.LinkedInURL = append(csvDataStruct.LinkedInURL, value)
			case "Email":
				csvDataStruct.Emails = append(csvDataStruct.Emails, value)
			}
		}
	}

	var resp []models.Payload
	var data models.Payload
	if len(csvDataStruct.Liid) > 0 { //unique identifier
		for idx := 0; idx < len(csvDataStruct.Liid); idx++ {
			if csvDataStruct.Liid[idx] == "" {
				if len(csvDataStruct.Emails) > 0 {
					if csvDataStruct.Emails[idx] != "" {
						fmt.Println("Emails")
						data, err = l.service.ScanDB(ctx, csvDataStruct.Emails[idx], "email")
						if err != nil {
							return CSVFileData{}, err
						}
						resp = append(resp, models.Payload{
							Emails:    data.Emails,
							Telephone: data.Telephone,
						})
					}
				}
				// phone hua for that idx
				if len(csvDataStruct.PhoneNumbers) > 0 {
					if csvDataStruct.PhoneNumbers[idx] != "" {
						data, err = l.service.ScanDB(ctx, csvDataStruct.PhoneNumbers[idx], "phone")
						if err != nil {
							return CSVFileData{}, err
						}
						resp = append(resp, models.Payload{
							Emails:    data.Emails,
							Telephone: data.Telephone,
						})
					}
				} else {
					continue
				}
			}

			data, err = l.service.ScanDB(ctx, csvDataStruct.Liid[idx], "liid")
			if err != nil {
				return CSVFileData{}, err
			}
			resp = append(resp, models.Payload{
				Emails:    data.Emails,
				Telephone: data.Telephone,
			})
		}
		for _, v := range resp {
			for _, email := range v.Emails {
				// check if email has gmail or yahoo or outlook or hotmail or icloud or aol or protonmail or zoho
				if strings.Contains(email, "gmail") || strings.Contains(email, "yahoo") || strings.Contains(email, "outlook") || strings.Contains(email, "hotmail") || strings.Contains(email, "icloud") || strings.Contains(email, "aol") || strings.Contains(email, "protonmail") || strings.Contains(email, "zoho") {
					csvDataStruct.PersonalEmails = append(csvDataStruct.PersonalEmails, email)
				} else {
					csvDataStruct.ProfessionalEmails = append(csvDataStruct.ProfessionalEmails, email)
				}
			}
		}

		return csvDataStruct, nil

	}

	if len(csvDataStruct.Emails) > 0 {
		for idx := 0; idx < len(csvDataStruct.Emails); idx++ {
			if csvDataStruct.Emails[idx] == "" {
				if len(csvDataStruct.PhoneNumbers) > 0 {
					if csvDataStruct.PhoneNumbers[idx] != "" {
						data, err = l.service.ScanDB(ctx, csvDataStruct.PhoneNumbers[idx], "phone")
						if err != nil {
							return CSVFileData{}, err
						}
						resp = append(resp, models.Payload{
							Emails:    data.Emails,
							Telephone: data.Telephone,
						})
					}
				} else {
					continue
				}
			}
			data, err := l.service.ScanDB(ctx, csvDataStruct.Emails[idx], "email")
			if err != nil {
				return CSVFileData{}, err
			}
			resp = append(resp, models.Payload{
				Emails:    data.Emails,
				Telephone: data.Telephone,
			})
		}

		for _, v := range resp {
			for _, email := range v.Emails{
				if strings.Contains(email, "gmail") || strings.Contains(email, "yahoo") || strings.Contains(email, "outlook") || strings.Contains(email, "hotmail") || strings.Contains(email, "icloud") || strings.Contains(email, "aol") || strings.Contains(email, "protonmail") || strings.Contains(email, "zoho") {
					csvDataStruct.PersonalEmails = append(csvDataStruct.PersonalEmails, email)
				} else {
					csvDataStruct.ProfessionalEmails = append(csvDataStruct.ProfessionalEmails, email)
				}
			}
		}

		return csvDataStruct, nil
	}

	if len(csvDataStruct.PhoneNumbers) > 0 {
		for idx := 0; idx < len(csvDataStruct.PhoneNumbers); idx++ {
			data, err := l.service.ScanDB(ctx, csvDataStruct.PhoneNumbers[idx], "phone")
			if err != nil {
				return CSVFileData{}, err
			}
			resp = append(resp, models.Payload{
				Emails:    data.Emails,
				Telephone: data.Telephone,
			})
		}

		for _, v := range resp {
			for _, email := range v.Emails {
				if strings.Contains(email, "gmail") || strings.Contains(email, "yahoo") || strings.Contains(email, "outlook") || strings.Contains(email, "hotmail") || strings.Contains(email, "icloud") || strings.Contains(email, "aol") || strings.Contains(email, "protonmail") || strings.Contains(email, "zoho") {
					csvDataStruct.PersonalEmails = append(csvDataStruct.PersonalEmails, email)
				} else {
					csvDataStruct.ProfessionalEmails = append(csvDataStruct.ProfessionalEmails, email)
				}
			}
		}

		return csvDataStruct, nil
	}

	return csvDataStruct, nil
}
