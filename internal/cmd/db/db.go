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

// FindInPep1 finds a document in the "pep1" collection based on LinkedIn identifier
func (s *Storage) FindInPep1(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error) {
	// pep1Data := models.Pep1Response{}
	// useData := models.UseResponse{}
	resp := models.Payload{}

	// collection := s.client.Database(s.dbName).Collection("pep1")
	// filter := bson.M{"liid": linkedInID}
	// projection := bson.D{{"e", 1}, {"t", 1}}
	// err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&pep1Data)
	// resp = models.Payload(pep1Data)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		collection = s.client.Database(s.dbName).Collection("use")
	// 		filter = bson.M{"linkedin_username": linkedInID}
	// 		projection = bson.D{{"emails", 1}, {"phone_numbers", 1}}
	// 		err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&useData)
	// 		resp.Emails = useData.Emails
	// 		resp.Telephone = useData.Telephone
	// 		if err != nil {
	// 			if err == mongo.ErrNoDocuments {
	// 				nullcounter++
	// 				fmt.Println(nullcounter)
	// 			}
	// 		}
	// 	}
	// }

	queries := []Query{
		// {Collection: "ap1", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "pep1", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "use", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "use1", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "use2", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
	}

	var wg sync.WaitGroup
	resultChan := make(chan models.Pep1Response)
	errChan := make(chan error)
	doneChan := make(chan bool)

	for _, q := range queries {
		wg.Add(1)

		go func(query Query) {
			defer wg.Done()
			c := s.client.Database(s.dbName).Collection(query.Collection)
			var result models.Pep1Response
			err := c.FindOne(ctx, query.Filter, options.FindOne().SetProjection(query.Projection)).Decode(&result)

			// If no data found, increment counter
			if err == mongo.ErrNoDocuments {
				mu.Lock()
				noDataCounter++
				// fmt.Println(noDataCounter)
				mu.Unlock()
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
		// wg.Done()
	}()

	func() {
		for {
			select {
			case err := <-errChan:
				log.Fatal("Error finding document: ", err)
			case result := <-resultChan:
				resp = models.Payload(result)
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
