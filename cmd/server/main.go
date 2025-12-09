// @title Employee Management API
// @version 1.0
// @description REST API for Employee Management
// @host localhost:8080	
// @BasePath /api/v1
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/employee_management_system/internal/app"
	"github.com/mohamedfawas/employee_management_system/internal/config"
)

func main() {

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = ""
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("[SERVER] Failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv, err := app.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("[SERVER] Failed to create server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("[SERVER] Failed to start HTTP server: %v", err)
				stop() // stops listening for signals, and exits the program
			}
		}
	}()

	log.Printf("[SERVER] HTTP server started on port %s", cfg.HTTP.Port)

	// Wait for shutdown signal
	<-ctx.Done()

	log.Println("[SERVER] Shutting down server...")
	if err := srv.Stop(context.Background()); err != nil {
		log.Printf("[SERVER] Error during server shutdown: %v", err)
	}
	log.Println("[SERVER] Server stopped")
}
