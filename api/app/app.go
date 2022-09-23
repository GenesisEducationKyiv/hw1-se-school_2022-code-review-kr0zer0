package app

import (
	"api/config"
	"api/internal/controllers/http"
	"api/internal/infrastructure/cryptoProviders"
	"api/internal/infrastructure/mailers"
	"api/internal/infrastructure/repository/file"
	"api/internal/service"
	"github.com/mailjet/mailjet-apiv3-go"
)

func Run() error {
	cfg := config.GetConfig()
	mailjetClient := mailjet.NewMailjetClient(cfg.EmailSending.PublicKey, cfg.EmailSending.PrivateKey)
	mailer := mailers.NewMailjetMailer(cfg, mailjetClient)

	coinMarketCapProviderCreator := crypto_providers.NewCoinMarketCapProviderCreator(cfg)
	binanceProviderCreator := crypto_providers.NewBinanceProviderCreator(cfg)
	coinAPIProviderCreator := crypto_providers.NewCoinAPIProviderCreator(cfg)
	coinbaseProviderCreator := crypto_providers.NewCoinbaseProviderCreator(cfg)

	coinMarketCapProvider := coinMarketCapProviderCreator.CreateCryptoProvider()
	binanceProvider := binanceProviderCreator.CreateCryptoProvider()
	coinAPIProvider := coinAPIProviderCreator.CreateCryptoProvider()
	coinbaseProvider := coinbaseProviderCreator.CreateCryptoProvider()

	coinMarketCapChain := service.NewBaseCryptoChain(service.NewLoggingCryptoProvider(coinMarketCapProvider))
	binanceChain := service.NewBaseCryptoChain(service.NewLoggingCryptoProvider(binanceProvider))
	coinAPIChain := service.NewBaseCryptoChain(service.NewLoggingCryptoProvider(coinAPIProvider))
	coinbaseChain := service.NewBaseCryptoChain(service.NewLoggingCryptoProvider(coinbaseProvider))

	coinMarketCapChain.SetNext(binanceChain)
	binanceChain.SetNext(coinAPIChain)
	coinAPIChain.SetNext(coinbaseChain)

	repos := file.NewRepository(cfg.Database.FilePath)
	services := service.NewService(repos, coinMarketCapChain, mailer, cfg)
	handlers := http.NewHandler(services)
	router := handlers.InitRouter()
	err := router.Run(cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
