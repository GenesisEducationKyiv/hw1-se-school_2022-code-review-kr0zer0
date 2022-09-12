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

	var cryptoProviderCreator service.CryptoProviderCreator
	switch cfg.CryptoProviders.CryptoProvider {
	case "Binance":
		cryptoProviderCreator = service.NewBinanceProviderCreator(cfg)
	case "CoinMarketCap":
		cryptoProviderCreator = service.NewCoinMarketCapProviderCreator(cfg)
	}

	cryptoProvider := cryptoProviderCreator.CreateCryptoProvider()

	repos := repository.NewRepository(cfg.Database.FilePath, cfg, mailjetClient)
	services := service.NewService(repos, cryptoProvider)
	handlers := handler.NewHandler(services)
	router := handlers.InitRouter()
	err := router.Run(cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
