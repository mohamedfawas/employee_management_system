package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/mohamedfawas/employee_management_system/internal/config"
	redisClient "github.com/mohamedfawas/employee_management_system/pkg/cache"
	"github.com/mohamedfawas/employee_management_system/pkg/constants"
	postgresClient "github.com/mohamedfawas/employee_management_system/pkg/database/postgres"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server

	// Database clients
	postgresClient *postgresClient.Client
	redisClient    *redisClient.Client
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	server := &Server{
		config: cfg,
	}

	// Initialize database clients first
	if err := server.initClients(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize clients: %w", err)
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Set production mode
	if cfg.Environment == constants.EnvProduction {
		e.Debug = false
	}

	// Setup middleware
	setupMiddleware(e)

	// Setup routes
	setupRoutes(e, server)

	// Create HTTP server
	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:      e,
		ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTP.IdleTimeout) * time.Second,
	}

	return server, nil
}

// initClients initializes PostgreSQL and Redis clients
func (s *Server) initClients(ctx context.Context) error {
	// Initialize PostgreSQL client
	postgresCfg := postgresClient.Config{
		Host:     s.config.Postgres.Host,
		Port:     s.config.Postgres.Port,
		User:     s.config.Postgres.User,
		Password: s.config.Postgres.Password,
		DBName:   s.config.Postgres.DBName,
		SSLMode:  s.config.Postgres.SSLMode,
		TimeZone: s.config.Postgres.TimeZone,
	}

	pgClient, err := postgresClient.NewClient(postgresCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize PostgreSQL client: %w", err)
	}
	s.postgresClient = pgClient

	// Initialize Redis client
	redisCfg := redisClient.Config{
		Host:     s.config.Redis.Host,
		Port:     s.config.Redis.Port,
		Password: s.config.Redis.Password,
		DB:       s.config.Redis.DB,
	}

	rdClient, err := redisClient.NewClient(redisCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis client: %w", err)
	}
	s.redisClient = rdClient

	return nil
}

// setupMiddleware configures Echo middleware
func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
}

// setupRoutes configures all HTTP routes
func setupRoutes(e *echo.Echo, s *Server) {
	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		if !s.isReady() {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "unhealthy",
				"error":  "dependencies not ready",
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// API v1 routes
	v1 := e.Group("/api/v1")
	{
		// Add your routes here
		_ = v1 // placeholder to avoid unused variable
	}
}

// isReady checks if all dependencies are initialized
func (s *Server) isReady() bool {
	return s.postgresClient != nil && s.redisClient != nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	// Close database connections
	if s.postgresClient != nil {
		s.postgresClient.Close()
	}

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			return fmt.Errorf("failed to close Redis client: %w", err)
		}
	}

	// Shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(shutdownCtx)
}
