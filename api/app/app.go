package app

import (
	"api/config"
	"api/internal/handler"
	"api/internal/repository"
	"api/internal/service"

	"github.com/mailjet/mailjet-apiv3-go"
)

func Run() error {
	cfg := config.GetConfig()
	mailjetClient := mailjet.NewMailjetClient(cfg.EmailSending.PublicKey, cfg.EmailSending.PrivateKey)

	coinMarketCapProvider := service.NewCoinMarketCapProvider(cfg)

	repos := repository.NewRepository(cfg.Database.FilePath, cfg, mailjetClient)
	services := service.NewService(repos, coinMarketCapProvider)
	handlers := handler.NewHandler(services)
	router := handlers.InitRouter()
	err := router.Run(cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
