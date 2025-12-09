package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mohamedfawas/employee_management_system/docs"
	cacheadapter "github.com/mohamedfawas/employee_management_system/internal/adapter/cache"
	postgresAdapter "github.com/mohamedfawas/employee_management_system/internal/adapter/db"
	"github.com/mohamedfawas/employee_management_system/internal/config"
	httpRouter "github.com/mohamedfawas/employee_management_system/internal/delivery/http"
	customMiddleware "github.com/mohamedfawas/employee_management_system/internal/delivery/http/middleware"
	v1 "github.com/mohamedfawas/employee_management_system/internal/delivery/http/v1"
	employeeUsecase "github.com/mohamedfawas/employee_management_system/internal/usecase"
	redisClient "github.com/mohamedfawas/employee_management_system/pkg/cache"
	"github.com/mohamedfawas/employee_management_system/pkg/constants"
	postgresClient "github.com/mohamedfawas/employee_management_system/pkg/database/postgres"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server

	postgresClient *postgresClient.Client
	redisClient    *redisClient.Client
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	server := &Server{
		config: cfg,
	}

	if err := server.initClients(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize clients: %w", err)
	}

	e := echo.New()
	e.HideBanner = true

	if cfg.Environment == constants.EnvProduction {
		e.Debug = false
	}

	setupMiddleware(e)

	//  Swagger UI endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	employeeRepo := postgresAdapter.NewEmployeeRepository(server.postgresClient.Pool)
	redisAdapter := cacheadapter.NewRedisAdapter(server.redisClient)
	employeeUsecase := employeeUsecase.NewEmployeeUsecase(employeeRepo, redisAdapter)
	employeeHandler := v1.NewEmployeeHandler(employeeUsecase)
	httpRouter.RegisterRoutes(e, employeeHandler)

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:      e,
		ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTP.IdleTimeout) * time.Second,
	}

	return server, nil
}

func (s *Server) initClients(ctx context.Context) error {
	postgresCfg := postgresClient.Config{
		Host:     s.config.Postgres.Host,
		Port:     s.config.Postgres.Port,
		User:     s.config.Postgres.User,
		Password: s.config.Postgres.Password,
		DBName:   s.config.Postgres.DBName,
		SSLMode:  s.config.Postgres.SSLMode,
	}

	pgClient, err := postgresClient.NewClient(ctx, postgresCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize PostgreSQL client: %w", err)
	}
	s.postgresClient = pgClient

	redisCfg := redisClient.Config{
		Host:     s.config.Redis.Host,
		Port:     s.config.Redis.Port,
		Password: s.config.Redis.Password,
		DB:       s.config.Redis.DB,
	}

	rdClient, err := redisClient.NewClient(ctx, redisCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis client: %w", err)
	}
	s.redisClient = rdClient

	return nil
}

func setupMiddleware(e *echo.Echo) {
	e.Use(customMiddleware.RequestIDMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.postgresClient != nil {
		s.postgresClient.Close()
	}

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			return fmt.Errorf("failed to close Redis client: %w", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(shutdownCtx)
}
