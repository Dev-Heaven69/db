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

func (s *Storage) FindInPep1(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error) {
	// pep1Data := models.Pep1Response{}
	// useData := models.UseResponse{}
	resp := models.Payload{}

	queries := []Query{
		{Collection: "ap2", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "use", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "use1", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
		{Collection: "use2", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "pep2personal", Filter: bson.D{{"liid", linkedInID}}, Projection: bson.D{{"e", 1}, {"t", 1}}},
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
				mu.Unlock()
				return
			}

			if err != nil {
                fmt.Println("error occured in collection: ",query.Collection, " while finding document: ", err)
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

func (s *Storage) GetPersonalEmail(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error) {
	emailRegex := bson.D{{"e", bson.D{{"$regex", `@(gmail\.com|hotmail\.me|yahoo\.in)$`}}}}

    queries := []Query{
        {Collection: "pep2personal", Filter: bson.D{{"liid", linkedInID}, emailRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "ap2", Filter: bson.D{{"liid", linkedInID}, emailRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "use", Filter: bson.D{{"liid", linkedInID}, emailRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "use1", Filter: bson.D{{"liid", linkedInID}, emailRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "use2", Filter: bson.D{{"liid", linkedInID}, emailRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
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

            if err == mongo.ErrNoDocuments {
                mu.Lock()
                noDataCounter++
                mu.Unlock()
                return
            }

            if err != nil {
                fmt.Println("error occured in collection: ",query.Collection, " while finding document: ", err)
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
            return models.Payload(result), nil
        case <-doneChan:
            if len(resultChan) == 0 {
                return models.Payload{}, fmt.Errorf("no personal emails found")
            }
            continue // Continue waiting until a result is available or all operations are completed
        }
    }
}

func (s *Storage) GetProfessionalEmails(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error) {
    // Using $not with $regex to exclude specified email domains
    professionalEmailsRegex := bson.D{{"e", bson.D{{"$not", bson.D{{"$regex", `@(gmail\.com|hotmail\.me|yahoo\.in)$`}}}}}}

    queries := []Query{
        {Collection: "ap2", Filter: bson.D{{"liid", linkedInID}, professionalEmailsRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "use", Filter: bson.D{{"liid", linkedInID}, professionalEmailsRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "use1", Filter: bson.D{{"liid", linkedInID}, professionalEmailsRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "use2", Filter: bson.D{{"liid", linkedInID}, professionalEmailsRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
        {Collection: "pep2personal", Filter: bson.D{{"liid", linkedInID}, professionalEmailsRegex[0]}, Projection: bson.D{{"e", 1}, {"t", 1}}},
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
            return models.Payload(result), nil
        case <-doneChan:
            if len(resultChan) == 0 {
                return models.Payload{}, fmt.Errorf("no emails found that do not match specified domains")
            }
            continue
        }
    }
}
