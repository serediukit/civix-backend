package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func (s *Store) Close() {
	s.db.Close()
}

func NewDB(ctx context.Context, config *DatabaseConfig) (*Store, error) {
	connStr := fmt.Sprintf(
		"db://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
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

func (s *Store) GetDB() *pgxpool.Pool {
	return s.db
}
