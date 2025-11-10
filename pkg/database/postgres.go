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
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
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
