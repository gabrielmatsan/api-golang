package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrator gerencia as migrations do banco
type Migrator struct {
	db *sql.DB
	m  *migrate.Migrate
}

// NewMigrator cria uma nova instância do migrator
func NewMigrator(db *sql.DB) (*Migrator, error) {
	// 1. Configurar source (arquivos embedded)
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	// 2. Configurar database driver
	databaseDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create database driver: %w", err)
	}

	// 3. Criar instância do migrate
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", databaseDriver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &Migrator{db: db, m: m}, nil
}

// Up executa todas as migrations pending
func (migrator *Migrator) Up() error {
	err := migrator.m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// Down executa rollback de todas as migrations
func (migrator *Migrator) Down() error {
	err := migrator.m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}
	return nil
}

// Steps executa um número específico de migrations
func (migrator *Migrator) Steps(n int) error {
	err := migrator.m.Steps(n)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run %d steps: %w", n, err)
	}
	return nil
}

// Version retorna a versão atual do banco
func (migrator *Migrator) Version() (uint, bool, error) {
	return migrator.m.Version()
}

// Close fecha a conexão do migrator
func (migrator *Migrator) Close() error {
	sourceErr, dbErr := migrator.m.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}

// ========================================
// COMANDOS CLI ÚTEIS
// ========================================

// Criar nova migration:
// migrate create -ext sql -dir migrations -seq create_users

// Executar migrations:
// migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up

// Rollback:
// migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down 1

// Ver versão atual:
// migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" version

// Forçar versão (em caso de erro):
// migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" force 1
