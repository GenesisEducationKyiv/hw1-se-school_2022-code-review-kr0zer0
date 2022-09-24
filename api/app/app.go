package app

import (
	"api/config"
	"api/internal/controllers/http"
	"api/internal/infrastructure/cryptoProviders"
	"api/internal/infrastructure/mailers"
	"api/internal/infrastructure/repository/file"
	"api/internal/service"
	"api/internal/service/crypto"
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

	coinMarketCapChain := crypto.NewBaseCryptoChain(crypto.NewLoggingCryptoProvider(coinMarketCapProvider))
	binanceChain := crypto.NewBaseCryptoChain(crypto.NewLoggingCryptoProvider(binanceProvider))
	coinAPIChain := crypto.NewBaseCryptoChain(crypto.NewLoggingCryptoProvider(coinAPIProvider))
	coinbaseChain := crypto.NewBaseCryptoChain(crypto.NewLoggingCryptoProvider(coinbaseProvider))

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
