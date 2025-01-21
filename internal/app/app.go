package app

import (
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/config"
	"github.com/Zrossiz/gophkeeper/internal/service"
	"github.com/Zrossiz/gophkeeper/internal/storage/postgres"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/handler"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/router"
	"github.com/Zrossiz/gophkeeper/pkg/logger"
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

	DBstore := postgres.New()
	serv := service.New()
	handler := handler.New()
	router := router.New()
}
