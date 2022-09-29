package app

import (
	"api/config"
	"api/internal/controllers/http"
	crypto_providers "api/internal/infrastructure/cryptoProviders"
	"api/internal/infrastructure/mailers"
	"api/internal/infrastructure/repository/filestorage"
	"api/internal/usecases"
	"api/internal/usecases/details"
	"api/internal/usecases/usecases_contracts"

	"github.com/mailjet/mailjet-apiv3-go"
)

func Run() error {
	cfg := config.GetConfig()

	mailer := initMailer(cfg)

	cryptoProvidersChain := initCryptoProvidersChain(cfg)

	repos := initRepos(cfg.Database.FilePath)
	useCases := initUseCases(repos, cryptoProvidersChain, mailer, cfg)
	handlers := http.NewHandler(useCases)

	router := handlers.InitRouter()

	err := router.Run(cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}

func initMailer(cfg *config.Config) usecases_contracts.Mailer {
	mailjetClient := mailjet.NewMailjetClient(cfg.EmailSending.PublicKey, cfg.EmailSending.PrivateKey)
	mailer := mailers.NewMailjetMailer(cfg, mailjetClient)

	return mailer
}

func initCryptoProvidersChain(cfg *config.Config) details.CryptoChain {
	coinMarketCapProviderCreator := crypto_providers.NewCoinMarketCapProviderCreator(cfg)
	binanceProviderCreator := crypto_providers.NewBinanceProviderCreator(cfg)
	coinAPIProviderCreator := crypto_providers.NewCoinAPIProviderCreator(cfg)
	coinbaseProviderCreator := crypto_providers.NewCoinbaseProviderCreator(cfg)

	coinMarketCapProvider := coinMarketCapProviderCreator.CreateCryptoProvider()
	binanceProvider := binanceProviderCreator.CreateCryptoProvider()
	coinAPIProvider := coinAPIProviderCreator.CreateCryptoProvider()
	coinbaseProvider := coinbaseProviderCreator.CreateCryptoProvider()

	coinMarketCapChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(coinMarketCapProvider))
	binanceChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(binanceProvider))
	coinAPIChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(coinAPIProvider))
	coinbaseChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(coinbaseProvider))

	coinMarketCapChain.SetNext(binanceChain)
	binanceChain.SetNext(coinAPIChain)
	coinAPIChain.SetNext(coinbaseChain)

	return coinMarketCapChain
}

func initRepos(filePath string) *usecases_contracts.Repository {
	emailSub := filestorage.NewEmailSubscriptionRepository(filePath)

	return filestorage.NewRepository(emailSub)
}

func initUseCases(
	repositories *usecases_contracts.Repository,
	cryptoChain details.CryptoChain,
	mailer usecases_contracts.Mailer,
	cfg *config.Config,
) *usecases.UseCases {
	getRate := details.NewCachedRateGetter(usecases.NewGetRateUseCase(cryptoChain), cfg.Cache.RateCacheTTL)
	sendEmails := usecases.NewSendEmailsUseCase(repositories.EmailSubscriptionRepo, mailer, getRate)
	subscribeEmails := usecases.NewSubscribeEmailUseCase(repositories.EmailSubscriptionRepo)

	return usecases.NewUseCases(getRate, sendEmails, subscribeEmails)
}
