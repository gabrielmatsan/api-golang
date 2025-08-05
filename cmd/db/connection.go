package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ===== DB Connection Config ===== //
type DBConnectionConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
}

func LoadDatabaseConfig() *DBConnectionConfig {
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("AVISO: Erro ao converter a porta do banco de dados '%s': %v. Usando porta padrão 5432.", portStr, err)
		port = 5432
	}

	return &DBConnectionConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         port,
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Database:     os.Getenv("DB_NAME"),
		SSLMode:      os.Getenv("DB_SSLMODE"),
		MaxOpenConns: 25,              // Número máximo de conexões abertas com o banco de dados.
		MaxIdleConns: 5,               // Número máximo de conexões ociosas no pool.
		MaxLifetime:  5 * time.Minute, // Tempo máximo que uma conexão pode ser reutilizada.
		MaxIdleTime:  1 * time.Minute, // Tempo máximo que uma conexão pode ficar ociosa antes de ser fechada.
	}
}

func (c *DBConnectionConfig) DatabaseURL() string {
	// Formato DSN comum para PostgreSQL:
	// "host=localhost port=5432 user=user password=password dbname=database sslmode=disable"
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// NewDbConnection estabelece uma nova conexão com o banco de dados e retorna uma instância sqlx.DB.
func NewDbConnection(config DBConnectionConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", config.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)
	db.SetConnMaxIdleTime(config.MaxIdleTime)

	// Ping para verificar a conexão imediatamente
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao fazer ping no banco de dados: %w", err)
	}

	return db, nil
}

func GetDB() *sqlx.DB {
	db, err := NewDbConnection(*LoadDatabaseConfig())
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	return db
}
