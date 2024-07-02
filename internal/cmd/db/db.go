package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/DevHeaven/db/domain/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client  *mongo.Client
	dbName  string
	timeout time.Duration
}

var noDataCounter int
var mu sync.Mutex

type Query struct {
	Collection string
	Filter     bson.D
	Projection bson.D
}

// Create a Mongo Client
func newMongoClient(dbUri string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to MongoDB")
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "failed to ping MongoDB")
	}
	fmt.Println("Connected to MongoDB!")

	return client, nil
}

// NewMongoRepository creates a new MongoDB repository
func NewMongoRepository(dbUri, dbName string, timeout int) (Storage, error) {
	mongoClient, err := newMongoClient(dbUri, timeout)
	if err != nil {
		return Storage{}, errors.Wrap(err, "failed to create MongoDB client")
	}

	repo := Storage{
		client:  mongoClient,
		dbName:  dbName,
		timeout: time.Duration(timeout) * time.Second,
	}

	return repo, nil
}

var nhi int = 0

func (s *Storage) ScanDB(ctx context.Context, uniqueID string, idType string,wantedFields map[string]bool) (models.Payload, error) {
	resp := models.Payload{}
	var queries []Query
	var projections bson.D

	for k, v := range wantedFields {
		if v{
			fmt.Println(k, " : boobs")
			if k == "PersonalEmail" || k == "ProfessionalEmail" || k == "e" {
				if projections.Map()["e"] == 1 {
					continue
				}
				projections = append(projections, bson.E{"e", 1})
			}else{
				projections = append(projections, bson.E{k, 1})
			}
		}
	}

	if wantedFields["Organization Name"] || wantedFields["Organization Domain"] { // ap2 and ap3
		if !wantedFields["linkedin"] {	
			if idType == "liid" {
				queries = []Query{
					{Collection: "ap3", Filter: bson.D{{idType,uniqueID}}, Projection: projections},
					{Collection: "ap2", Filter: bson.D{{idType,uniqueID}}, Projection: projections},
				}
			}
		}
		if !wantedFields["e"] {
			if idType == "email" {
				queries = []Query{
					{Collection: "ap3", Filter: bson.D{{"e",uniqueID}}, Projection: projections},
					{Collection: "ap2", Filter: bson.D{{"e",uniqueID}}, Projection: projections},
				}
			}
		}
		if !wantedFields["t"] {
			if idType == "phone" {
				queries = []Query{
					{Collection: "ap3", Filter: bson.D{{"t",uniqueID}}, Projection: projections},
					{Collection: "ap2", Filter: bson.D{{"t",uniqueID}}, Projection: projections},
				}
			}
		}
	}else{ 	// pe1 and ap2
		if !wantedFields["linkedin"] {
			if idType == "liid" {
				queries = []Query{
					{Collection: "pe1", Filter: bson.D{{idType,uniqueID}}, Projection: projections},
					{Collection: "ap2", Filter: bson.D{{idType,uniqueID}}, Projection: projections},
					{Collection: "ap3", Filter: bson.D{{idType,uniqueID}}, Projection: projections},
				}
			}
		}
		if !wantedFields["e"] {
			if idType == "email" {
				queries = []Query{
					{Collection: "pe1", Filter: bson.D{{"e",uniqueID}}, Projection: projections},
					{Collection: "ap2", Filter: bson.D{{"e",uniqueID}}, Projection: projections},
					{Collection: "ap3", Filter: bson.D{{"e",uniqueID}}, Projection: projections},
				}
			}
		}
		if !wantedFields["t"] {
			if idType == "phone" {
				queries = []Query{
					{Collection: "pe1", Filter: bson.D{{"t",uniqueID}}, Projection: projections},
					{Collection: "ap2", Filter: bson.D{{"t",uniqueID}}, Projection: projections},
					{Collection: "ap3", Filter: bson.D{{"t",uniqueID}}, Projection: projections},
				}
			}
		}
	}

	var wg sync.WaitGroup
	resultChan := make(chan models.DbResponse)
	errChan := make(chan error)
	doneChan := make(chan bool)

	for _, q := range queries {
		wg.Add(1)

		go func(query Query) {
			defer wg.Done()
			c := s.client.Database(s.dbName).Collection(query.Collection)
			var result models.DbResponse
			err := c.FindOne(ctx, query.Filter, options.FindOne().SetProjection(query.Projection)).Decode(&result)

			// If no data found, increment counter
			if err == mongo.ErrNoDocuments {
				mu.Lock()
				noDataCounter++
				mu.Unlock()
				return
			}

			if err != nil {
				fmt.Println("error occured in collection: ", query.Collection, " while finding document: ", err , " ", query)
				errChan <- err
				return
			}

			resultChan <- result
			doneChan <- true
		}(q)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	func() {
		for {
			select {
			case err := <-errChan:
				log.Fatal("Error finding document: ", err)
			case result := <-resultChan:
				resp = models.Payload{
					Emails: result.Emails,
					Telephone: result.Telephone,
					OrganizationDomain: result.OrganizationDomain,
					OrganizationName: result.OrganizationName,
					LinkedInUrl: result.LinkedInUrl,
					FirstName: result.FirstName,
					LastName: result.LastName,
				}
				return
			case <-doneChan:
				nhi++
				fmt.Println(nhi)
				return
			}
		}
	}()

	return resp, nil
}


func (s *Storage) GetPersonalEmail(ctx context.Context, linkedInID string) (models.Payload, error) {
	emailRegex := bson.D{{"e", bson.D{{"$regex", `@(gmail\.com|hotmail\.me|yahoo\.in)$`}}}}

	queries := []Query{
		{Collection: "ap2", Filter: bson.D{{"liid", linkedInID}, emailRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "pe1", Filter: bson.D{{"liid", linkedInID}, emailRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
	}
	var wg sync.WaitGroup
	resultChan := make(chan models.DbResponse)
	errChan := make(chan error)
	doneChan := make(chan bool)

	for _, q := range queries {
		wg.Add(1)

		go func(query Query) {
			defer wg.Done()
			c := s.client.Database(s.dbName).Collection(query.Collection)
			var result models.DbResponse
			err := c.FindOne(ctx, query.Filter, options.FindOne().SetProjection(query.Projection)).Decode(&result)

			if err == mongo.ErrNoDocuments {
				mu.Lock()
				noDataCounter++
				mu.Unlock()
				return
			}

			if err != nil {
				fmt.Println("error occured in collection: ", query.Collection, " while finding document: ", err)
				errChan <- err
				return
			}

			resultChan <- result
			doneChan <- true
		}(q)
	}

	// Waiting for routines to finish
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	// Blocking main thread to handle results and errors
	for {
		select {
		case err := <-errChan:
			return models.Payload{}, err
		case result := <-resultChan:
			return models.Payload{
				Emails: result.Emails,
				Telephone: result.Telephone,
			}, nil
		case <-doneChan:
			if len(resultChan) == 0 {
				return models.Payload{}, fmt.Errorf("no personal emails found")
			}
			continue // Continue waiting until a result is available or all operations are completed
		}
	}
}

func (s *Storage) GetProfessionalEmails(ctx context.Context, linkedInID string) (models.Payload, error) {
	// Using $not with $regex to exclude specified email domains
	professionalEmailsRegex := bson.D{{"e", bson.D{{"$not", bson.D{{"$regex", `@(gmail\.com|hotmail\.me|yahoo\.in)$`}}}}}}

	queries := []Query{
		{Collection: "ap2", Filter: bson.D{{"liid", linkedInID}, professionalEmailsRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "pe1", Filter: bson.D{{"liid", linkedInID}, professionalEmailsRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
	}

	var wg sync.WaitGroup
	resultChan := make(chan models.DbResponse)
	errChan := make(chan error)
	doneChan := make(chan bool)

	for _, q := range queries {
		wg.Add(1)

		go func(query Query) {
			defer wg.Done()
			c := s.client.Database(s.dbName).Collection(query.Collection)
			var result models.DbResponse
			err := c.FindOne(ctx, query.Filter, options.FindOne().SetProjection(query.Projection)).Decode(&result)

			if err == mongo.ErrNoDocuments {
				return
			}

			if err != nil {
				errChan <- err
				return
			}

			resultChan <- result
			doneChan <- true
		}(q)
	}

	// Waiting for routines to finish
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	for {
		select {
		case err := <-errChan:
			return models.Payload{}, err
		case result := <-resultChan:
			return models.Payload{
				Emails: result.Emails,
				Telephone: result.Telephone,
			}, nil
		case <-doneChan:
			if len(resultChan) == 0 {
				return models.Payload{}, fmt.Errorf("no emails found that do not match specified domains")
			}
			continue
		}
	}
}
