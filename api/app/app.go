package app

import (
	"api/config"
	"api/internal/handler"
	"api/internal/repository"
	"api/internal/service"
)

func Run() error {
	cfg := config.GetConfig()

	repos := repository.NewRepository(cfg.Database.FilePath)
	services := service.NewService(repos, cfg)
	handlers := handler.NewHandler(services)
	err := handlers.InitRouter(cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
