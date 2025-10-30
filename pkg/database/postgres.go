package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfigGetter interface {
	GetDBConfig() *DatabaseConfig
}

type Store struct {
	db *pgxpool.Pool
}

func NewDB(ctx context.Context, config DBConfigGetter) (*Store, error) {
	dbConfig := config.GetDBConfig()

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.SSLMode,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("[PostgreSQL] Unable to connect: %v\n", err)
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("[PostgreSQL] Unable to ping: %v\n", err)
		return nil, err
	}

	return &Store{
		db: pool,
	}, nil
}
