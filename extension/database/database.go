package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseHost string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort string `env:"DATABASE_PORT" envDefault:"5433"`
	DatabaseUser string `env:"DATABASE_USER" envDefault:"postgres"`
	DatabasePass string `env:"DATABASE_PASS" envDefault:"postgres"`
	DatabaseName string `env:"DATABASE_NAME" envDefault:"post_confiavel"`
}

func NewDatabase() (*sql.DB, error) {
	_ = godotenv.Load()

	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("database config: %w", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabasePass,
		config.DatabaseName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("sql ping: %w", err)
	}

	log.Println("connected to PostgreSQL")

	return db, nil
}
