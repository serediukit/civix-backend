package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestContainers struct {
	PostgresContainer *postgres.PostgresContainer
	RedisContainer    *redis.RedisContainer
	PostgresURL       string
	RedisURL          string
}

func SetupTestContainers(ctx context.Context) (*TestContainers, error) {
	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgis/postgis:15-3.3",
		postgres.WithDatabase("civix_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	pgURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres connection string: %w", err)
	}

	// Start Redis container
	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	redisURL, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis connection string: %w", err)
	}

	return &TestContainers{
		PostgresContainer: pgContainer,
		RedisContainer:    redisContainer,
		PostgresURL:       pgURL,
		RedisURL:          redisURL,
	}, nil
}

func (tc *TestContainers) Teardown(ctx context.Context) error {
	if tc.PostgresContainer != nil {
		if err := tc.PostgresContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate postgres container: %w", err)
		}
	}
	if tc.RedisContainer != nil {
		if err := tc.RedisContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate redis container: %w", err)
		}
	}
	return nil
}
