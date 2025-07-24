package main

import (
	"log"
	"os"

	"github.com/gabrielmatsan/teste-api/cmd/db"
	"github.com/gabrielmatsan/teste-api/cmd/server"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("AVISO: Erro ao carregar arquivo .env: %v", err)
	}

	dbConfig := db.LoadDatabaseConfig()
	db, err := db.NewDbConnection(*dbConfig)
	if err != nil {
		panic(err)
	}

	server := server.NewServer(db, os.Getenv("API_PORT"))
	server.SetupMiddlewares()
	server.SetupRoutes()

	serverInstance := server.Start()

	if err := serverInstance; err != nil {
		panic(err)
	}
}
