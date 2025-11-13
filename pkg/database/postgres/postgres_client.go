package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds PostgreSQL connection parameters.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// Client wraps a pgxpool.Pool connection.
type Client struct {
	Pool *pgxpool.Pool
}

// NewClient initializes and returns a PostgreSQL client using pgxpool.
// It performs connection pooling, health checks, and timeouts.
func NewClient(cfg Config) (*Client, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&TimeZone=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	// Add connection timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL config: %w", err)
	}

	// Establish pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL pool: %w", err)
	}

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully")

	return &Client{Pool: pool}, nil
}

// Close gracefully shuts down the connection pool.
func (c *Client) Close() {
	if c.Pool != nil {
		c.Pool.Close()
		log.Println("🔒 PostgreSQL connection closed")
	}
}
