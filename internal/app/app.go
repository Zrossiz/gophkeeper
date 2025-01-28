package app

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/gophkeeper/internal/config"
	"github.com/Zrossiz/gophkeeper/internal/service"
	"github.com/Zrossiz/gophkeeper/internal/storage/postgres"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/handler"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/router"
	"github.com/Zrossiz/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

func Start() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println("init config error: ", err)
	}

	log, err := logger.New(cfg.LoggerLevel)
	if err != nil {
		fmt.Println("init logger error: ", err)
	}

	dbConn, err := postgres.Connect(cfg.DBURI)
	if err != nil {
		log.Error("db connection error", zap.Error(err))
	}
	defer dbConn.Close()

	dbStore := postgres.New(dbConn)
	serv := service.New(service.Storage{
		Card:     &dbStore.Card,
		User:     &dbStore.User,
		Binary:   &dbStore.Binary,
		LogoPass: &dbStore.LogoPass,
	}, *cfg, log)

	handler := handler.New(handler.Service{
		Card:     &serv.Card,
		User:     &serv.User,
		Binary:   &serv.Binary,
		LogoPass: &serv.LogoPass,
	}, log)

	router := router.New(router.Handler{
		Card:     handler.Card,
		User:     handler.User,
		Binary:   handler.Binary,
		LogoPass: handler.LogoPass,
	})

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	log.Sugar().Infof("Starting server on addr: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("start web server error", zap.Error(err))
	}
}
