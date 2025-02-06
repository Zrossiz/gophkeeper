// Package app initializes and starts the GophKeeper application.
// It sets up the configuration, logger, database connection, services, HTTP handlers, and middleware.
// The application provides an API for managing secret data such as cards, passwords, and notes.
//
// The package also includes Swagger documentation for the API.
package app

import (
	"fmt"
	"net/http"

	_ "github.com/Zrossiz/gophkeeper/docs"
	"github.com/Zrossiz/gophkeeper/internal/config"
	"github.com/Zrossiz/gophkeeper/internal/cryptox"
	"github.com/Zrossiz/gophkeeper/internal/service"
	"github.com/Zrossiz/gophkeeper/internal/storage/postgres"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/handler"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/middleware"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/router"
	"github.com/Zrossiz/gophkeeper/pkg/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Start initializes and starts the GophKeeper application.
// It performs the following steps:
//  1. Loads the application configuration.
//  2. Initializes the logger.
//  3. Establishes a connection to the PostgreSQL database.
//  4. Initializes the authentication middleware and cryptographic module.
//  5. Sets up the database storage, services, and HTTP handlers.
//  6. Configures the HTTP router with middleware and handlers.
//  7. Starts the HTTP server on the configured address.
//
// If any initialization step fails, the function logs the error and terminates the application.
//
// @title GophKeeper API
// @version 1.0
// @description API для управления секретными данными (карточки, пароли и т.д.).
// @host localhost:8080
// @BasePath /api
func Start() {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		fmt.Println("init config error: ", err)
	}

	// Initialize logger
	log, err := logger.New(cfg.LoggerLevel)
	if err != nil {
		fmt.Println("init logger error: ", err)
	}

	// Connect to the database
	dbConn, err := postgres.Connect(cfg.DBURI)
	if err != nil {
		log.Error("db connection error", zap.Error(err))
	}
	defer dbConn.Close()

	// Initialize middleware and cryptographic module
	authMiddleware := middleware.New(*cfg, log)
	cryptoModule := cryptox.NewCryproModule()

	// Initialize database storage
	dbStore := postgres.New(dbConn)

	// Initialize services
	serv := service.New(service.Storage{
		Card:     &dbStore.Card,
		User:     &dbStore.User,
		Binary:   &dbStore.Binary,
		LogoPass: &dbStore.LogoPass,
		Note:     &dbStore.Note,
	}, *cfg, cryptoModule, log)

	// Initialize HTTP handlers
	handler := handler.New(handler.Service{
		Card:     &serv.Card,
		User:     &serv.User,
		Binary:   &serv.Binary,
		LogoPass: &serv.LogoPass,
		Note:     &serv.Note,
	}, log)

	// Configure HTTP router
	router := router.New(router.Handler{
		Card:     &handler.Card,
		User:     &handler.User,
		Binary:   &handler.Binary,
		LogoPass: &handler.LogoPass,
		Note:     &handler.Note,
	}, authMiddleware)

	// Start HTTP server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	log.Sugar().Infof("Starting server on addr: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("start web server error", zap.Error(err))
	}
}
