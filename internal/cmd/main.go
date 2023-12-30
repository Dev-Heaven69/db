package main

import (
	"fmt"
	"os"

	"github.com/DevHeaven/db/internal/cmd/db"
	"github.com/DevHeaven/db/internal/cmd/server"
	"github.com/DevHeaven/db/internal/dbi"
	"github.com/DevHeaven/db/internal/logic"
	"github.com/DevHeaven/db/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	port := os.Getenv("port")
	dbUri := os.Getenv("dburi")
	dbName := os.Getenv("dbname")
	engine := gin.Default()
	database, err := db.NewMongoRepository(dbUri, dbName, 10)
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	Services := dbi.NewService(&database)
	Logic := logic.ProvideLogic(Services)
	Routes := router.ProvideRouter(Logic)
	server := server.NewServer(port, engine, Routes)
	server.Run()
}
