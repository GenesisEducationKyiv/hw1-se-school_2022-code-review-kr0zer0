package app

import (
	"api/internal/handler"
	"api/internal/repository"
	"api/internal/service"
)

func Run() error {
	repos := repository.NewRepository("./data/data.json")
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	err := handlers.InitRouter(":8000")
	if err != nil {
		return err
	}

	return nil
}
