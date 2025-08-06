package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gabrielmatsan/teste-api/cmd/db"
	"github.com/gabrielmatsan/teste-api/cmd/server"
	"github.com/joho/godotenv"
)


// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
func main() {

	fmt.Println("=== INICIO DO MAIN ===")
	fmt.Println("Iniciando servidor...")

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
