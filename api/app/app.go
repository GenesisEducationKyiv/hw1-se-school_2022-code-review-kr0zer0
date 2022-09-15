package app

import (
	"api/config"
	"api/internal/handler"
	"api/internal/repository"
	"api/internal/service"
	"api/internal/service/cryptoProviders"

	"github.com/mailjet/mailjet-apiv3-go"
)

func Run() error {
	cfg := config.GetConfig()
	mailjetClient := mailjet.NewMailjetClient(cfg.EmailSending.PublicKey, cfg.EmailSending.PrivateKey)

	coinMarketCapProviderCreator := crypto_providers.NewCoinMarketCapProviderCreator(cfg)
	binanceProviderCreator := crypto_providers.NewBinanceProviderCreator(cfg)
	coinAPIProviderCreator := crypto_providers.NewCoinAPIProviderCreator(cfg)
	coinbaseProviderCreator := crypto_providers.NewCoinbaseProviderCreator(cfg)

	coinMarketCapProvider := coinMarketCapProviderCreator.CreateCryptoProvider()
	binanceProvider := binanceProviderCreator.CreateCryptoProvider()
	coinAPIProvider := coinAPIProviderCreator.CreateCryptoProvider()
	coinbaseProvider := coinbaseProviderCreator.CreateCryptoProvider()

	coinMarketCapChain := service.NewBaseCryptoChain(coinMarketCapProvider)
	binanceChain := service.NewBaseCryptoChain(binanceProvider)
	coinAPIChain := service.NewBaseCryptoChain(coinAPIProvider)
	coinbaseChain := service.NewBaseCryptoChain(coinbaseProvider)

	coinMarketCapChain.SetNext(binanceChain)
	binanceChain.SetNext(coinAPIChain)
	coinAPIChain.SetNext(coinbaseChain)

	repos := repository.NewRepository(cfg.Database.FilePath, cfg, mailjetClient)
	services := service.NewService(repos, coinMarketCapChain, cfg)
	handlers := handler.NewHandler(services)
	router := handlers.InitRouter()
	err := router.Run(cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
