package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/serediukit/civix-backend/pkg/database"
	"github.com/serediukit/civix-backend/pkg/redis"
)

type TestDatabaseSetup struct {
	Store       *database.Store
	RedisClient *redis.CachedStore
	Containers  *TestContainers
}

func SetupTestDatabase(ctx context.Context) (*TestDatabaseSetup, error) {
	// Setup containers
	containers, err := SetupTestContainers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup containers: %w", err)
	}

	// Parse PostgreSQL DSN
	dsn := containers.PostgresURL

	// Initialize database store
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		containers.Teardown(ctx)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(ctx, db); err != nil {
		containers.Teardown(ctx)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create database store
	store, err := database.NewDBFromDSN(ctx, dsn)
	if err != nil {
		containers.Teardown(ctx)
		return nil, fmt.Errorf("failed to create database store: %w", err)
	}

	// Setup Redis client
	// Parse Redis URL to get host and port
	// URL format: redis://localhost:port
	redisURL := containers.RedisURL
	var redisHost, redisPort string

	// Simple parsing for redis://host:port format
	if len(redisURL) > 8 { // "redis://" is 8 chars
		hostPort := redisURL[8:] // Remove "redis://" prefix
		// Split by ':'
		parts := strings.Split(hostPort, ":")
		if len(parts) == 2 {
			redisHost = parts[0]
			redisPort = parts[1]
		} else {
			redisHost = "localhost"
			redisPort = "6379"
		}
	}

	redisClient, err := redis.NewRedis(&redis.RedisConfig{
		Host:         redisHost,
		Port:         redisPort,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
	})
	if err != nil {
		containers.Teardown(ctx)
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	return &TestDatabaseSetup{
		Store:       store,
		RedisClient: redisClient,
		Containers:  containers,
	}, nil
}

func (s *TestDatabaseSetup) Teardown(ctx context.Context) error {
	if s.Store != nil {
		s.Store.Close()
	}
	if s.RedisClient != nil {
		s.RedisClient.Close()
	}
	if s.Containers != nil {
		return s.Containers.Teardown(ctx)
	}
	return nil
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	// Get migration files
	migrationsPath := "../../migrations"
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort migration files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Execute migrations
	for _, fileName := range migrationFiles {
		filePath := filepath.Join(migrationsPath, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}
	}

	return nil
}
