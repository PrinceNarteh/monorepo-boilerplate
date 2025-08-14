// Package main is the entry point to the program
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PrinceNarteh/go-boilerplate/internal/config"
	// "github.com/PrinceNarteh/go-boilerplate/internal/database" // Uncomment when using database
	"github.com/PrinceNarteh/go-boilerplate/internal/logger"
	"github.com/PrinceNarteh/go-boilerplate/internal/middlewares"
	"github.com/PrinceNarteh/go-boilerplate/internal/routers"
	"github.com/PrinceNarteh/go-boilerplate/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger service
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	// Initialize logger
	appLogger := logger.NewLoggerWithService(cfg.Observability, loggerService)

	// Initialize database (uncomment when you have a database)
	// db, err := database.New(cfg, &appLogger, loggerService)
	// if err != nil {
	//     appLogger.Fatal().Err(err).Msg("Failed to initialize database")
	// }
	// defer db.Close()

	// Run migrations (uncomment when you have a database)
	// ctx := context.Background()
	// if err := database.Migrate(ctx, &appLogger, cfg); err != nil {
	//     appLogger.Fatal().Err(err).Msg("Failed to run migrations")
	// }

	// Initialize router
	router := routers.New(&appLogger)
	router.SetupRoutes()

	// Setup middleware chain
	middlewareChain := middlewares.Chain(
		middlewares.Recovery(&appLogger),
		middlewares.Logger(&appLogger),
		middlewares.CORS(cfg.Server.CORSAllowedOrigins),
	)

	// Apply middleware to router
	handler := middlewareChain(router)

	// Initialize and start server
	srv := server.New(cfg, handler, &appLogger)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			appLogger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info().Msg("Shutting down server...")

	// Give server 30 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		appLogger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	appLogger.Info().Msg("Server exited")
}
