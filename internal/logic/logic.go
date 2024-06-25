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
	// var apidata [][]string
	// var csvDataStruct CSVFileData
	for idx := 1; idx < len(fields); idx++ {
		// fmt.Println(idx)
		parts := strings.Split(fields[idx][4], "/")
		id := parts[len(parts)-1]
		data, err := l.service.ScanDB(ctx, id)
		if err != nil {
			return nil, err
		}
		// if data.Emails == nil {

		// 	suspect := []string{fields[idx][2], fields[idx][3], fields[idx][16]}
		// 	_ = append(apidata, suspect)
		// }
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

		data, err := l.service.ScanDB(ctx, id)
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
		data, err := l.service.ScanDB(ctx, id)
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
	data, err := l.service.ScanDB(ctx, liid)
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
		data, err := l.service.ScanDB(ctx, liid)
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

func (l Logic) Test(file *multipart.FileHeader, ctx *gin.Context) (string, error) {
	uploadPath := "./data/"
	filename := "req.csv"

	filepath := uploadPath + filename

	err := ctx.SaveUploadedFile(file, filepath)
	if err != nil {
		return "", err
	}

	csvFile, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	headers, err := reader.Read()
	if err != nil {
		return "", err
	}

	// for {
	// 	record, err := reader.Read()
	// 	if err != nil {
	// 		break
	// 	}

	// 	fields = append(fields, record)
	// }
	fields,err := reader.ReadAll()
	if err != nil {	
		return "", err
	}

	// var resp []models.Payload
	// var apidata [][]string
	// fmt.Println("fields",len(fields[0]))
	var csvDataStruct CSVFileData
	// for idx := 0; idx < len(fields[0]); idx++ {
	// 	var firstName []string
	// 	var lastName []string
	// 	var organizationDomain []string
	// 	var linkedInURL []string
	// 	var phoneNumbers []string
	// 	var liid []string
	// 	// fmt.Println(idx)
	// 	// fmt.Println("this is fields", fields[0][idx])
	// 	if fields[0][idx] == "First Name" {
	// 		for _, v := range fields {
	// 			// fmt.Println("this is k",k)
	// 			firstName = append(firstName, v[idx])
	// 			// fmt.Println("this is v",v[idx])
	// 		}
	// 		csvDataStruct.FirstName = firstName
	// 	}else{
	// 		firstName = make([]string, 0)
	// 		csvDataStruct.FirstName = firstName
	// 	}

	// 	if fields[0][idx] == "Last Name" {
			
	// 		for _, v := range fields {
	// 			// fmt.Println("this is k",k)
	// 			// fmt.Println("this is v",v[idx])
	// 			lastName = append(lastName, v[idx])
	// 		}
	// 		csvDataStruct.LastName = lastName
	// 	}else{
	// 		lastName = make([]string, 0)
	// 		csvDataStruct.LastName = lastName
	// 	}

	// 	if fields[0][idx] == "Organization Domain" {
			
	// 		for _, v := range fields {
	// 			// fmt.Println("this is k",k)
	// 			// fmt.Println("this is v",v[idx])
	// 			organizationDomain = append(organizationDomain, v[idx])
	// 		}
	// 		csvDataStruct.OrganizationDomain = organizationDomain
	// 	}else{
	// 		organizationDomain = make([]string, 0)
	// 		csvDataStruct.OrganizationDomain = organizationDomain
	// 	}

	// 	if fields[0][idx] == "Phone Number" {
	// 		for _, v := range fields {
	// 			// fmt.Println("this is k",k)
	// 			// fmt.Println("this is v",v[idx])
	// 			phoneNumbers = append(phoneNumbers, v[idx])
	// 		}
	// 		csvDataStruct.PhoneNumbers = phoneNumbers
	// 	}else{
	// 		phoneNumbers = make([]string, 0)
	// 		csvDataStruct.PhoneNumbers = phoneNumbers
	// 	}

	// 	if fields[0][idx] == "LinkedIn" {
	// 		for _, v := range fields {
	// 			linkedInURL = append(linkedInURL, v[idx])
	// 			parts := strings.Split(v[idx], "/")
	// 			liid = append(liid, parts[len(parts)-1])
	// 			// fmt.Println("this is k",k)
	// 			// fmt.Println("this is v",v[idx])
	// 		}
	// 		csvDataStruct.LinkedInURL = linkedInURL
	// 		csvDataStruct.Liid = liid
	// 	}else{
	// 		linkedInURL = make([]string, 0)
	// 		liid = make([]string, 0)
	// 		csvDataStruct.LinkedInURL = linkedInURL
	// 		csvDataStruct.Liid = liid
	// 	}
		// fmt.Println("this is fields",fields[idx][0])
		// parts := strings.Split(fields[idx][5], "/")
		// id := parts[len(parts)-1]
		// data, err := l.service.FindInPep1(ctx, id)
		// if err != nil {
		// 	return nil, err
		// }
	// }
	for _, record := range fields {
		for i , value := range record {
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
			}
		}
	}
	for i := 0; i < len(fields);i++ {
		if (len(csvDataStruct.FirstName) > 0) {
			fmt.Println("this is firstname",csvDataStruct.FirstName[i])
		}
		if (len(csvDataStruct.LastName) > 0) {
			fmt.Println("this is lastname",csvDataStruct.LastName[i])
		}
		if (len(csvDataStruct.Liid) > 0) {
			fmt.Println("this is liid",csvDataStruct.Liid[i])
		}
		if (len(csvDataStruct.LinkedInURL) > 0) {
			fmt.Println("this is linkedinURL",csvDataStruct.LinkedInURL[i])
		}
		if (len(csvDataStruct.OrganizationDomain) > 0) {
			fmt.Println("this is domain",csvDataStruct.OrganizationDomain[i])
		}
		if (len(csvDataStruct.PhoneNumbers) > 0) {
			fmt.Println("pn",csvDataStruct.PhoneNumbers[i])
		}
		// fmt.Println(i)
		// data,err := l.service.FindInPep1()

		// fmt.Println("this is firstname",csvDataStruct.FirstName[0])
		// fmt.Println("this is lastname",csvDataStruct.LastName[0])
		// fmt.Println("this is liid",csvDataStruct.Liid[0])
		// fmt.Println("this is linkedinURL",csvDataStruct.LinkedInURL[0])
		// fmt.Println("this is domain",csvDataStruct.OrganizationDomain[0])
		// fmt.Println("pn",csvDataStruct.PhoneNumbers[0])
		// fmt.Println("emaisl	",csvDataStruct.Emails[0])
	
	}
	// fmt.Println("=====================================")
	// fmt.Println("this is firstname",csvDataStruct.FirstName)
	// fmt.Println("=====================================")
	
	// fmt.Println("this is lastname",csvDataStruct.LastName)
	// fmt.Println("=====================================")

	// fmt.Println("this is liid",csvDataStruct.Liid)
	// fmt.Println("=====================================")
	// fmt.Println("this is linkedinURL",csvDataStruct.LinkedInURL)
	// fmt.Println("=====================================")
	// fmt.Println("this is domain",csvDataStruct.OrganizationDomain)
	// fmt.Println("=====================================")
	// fmt.Println("pn",csvDataStruct.PhoneNumbers)
	// fmt.Println("=====================================")
	// fmt.Println("emaisl	",csvDataStruct.Emails)
	// fmt.Println("=====================================")


	

	return "hello world", nil
}