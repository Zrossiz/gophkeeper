package app

import (
	"github.com/Zrossiz/gophkeeper/internal/config"
	"github.com/Zrossiz/gophkeeper/internal/service"
	"github.com/Zrossiz/gophkeeper/internal/storage/postgres"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/handler"
	"github.com/Zrossiz/gophkeeper/internal/transport/http/router"
)

func Start() {
	cfg := config.New()
	DBstore := postgres.New()
	serv := service.New()
	handler := handler.New()
	router := router.New()
}
