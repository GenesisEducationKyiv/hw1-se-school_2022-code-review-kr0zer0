package app

import (
	"api/config"
	"api/internal/controllers/http"
	"api/internal/infrastructure/brokers"
	crypto_providers "api/internal/infrastructure/cryptoProviders"
	"api/internal/infrastructure/mailers"
	"api/internal/infrastructure/repository/filestorage"
	"api/internal/usecases"
	"api/internal/usecases/details"
	"api/internal/usecases/usecases_contracts"
	"api/pkg/logging"
	"github.com/mailjet/mailjet-apiv3-go"
	amqp "github.com/rabbitmq/amqp091-go"

	//amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"io"
)

func Run() error {
	connection, err := amqp.Dial("amqp://test:test@rabbitmq:5672/")
	if err != nil {
		return err
	}
	defer connection.Close()

	brokerLogger := brokers.NewBrokerWriter(connection)
	logger := logging.Init([]io.Writer{brokerLogger})

	cfg := config.GetConfig()

	mailer := initMailer(cfg, logger)

	cryptoProvidersChain := initCryptoProvidersChain(cfg, logger)

	repos := initRepos(cfg.Database.FilePath, logger)
	useCases := initUseCases(repos, cryptoProvidersChain, mailer, cfg)
	handlers := http.NewHandler(useCases)

	router := handlers.InitRouter()

	err = router.Run(cfg.App.Port)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func initMailer(cfg *config.Config, logger *logrus.Logger) usecases_contracts.Mailer {
	mailjetClient := mailjet.NewMailjetClient(cfg.EmailSending.PublicKey, cfg.EmailSending.PrivateKey)
	mailer := mailers.NewMailjetMailer(cfg, mailjetClient, logger)

	return mailer
}

func initCryptoProvidersChain(cfg *config.Config, logger *logrus.Logger) details.CryptoChain {
	coinMarketCapProviderCreator := crypto_providers.NewCoinMarketCapProviderCreator(cfg)
	binanceProviderCreator := crypto_providers.NewBinanceProviderCreator(cfg)
	coinAPIProviderCreator := crypto_providers.NewCoinAPIProviderCreator(cfg)
	coinbaseProviderCreator := crypto_providers.NewCoinbaseProviderCreator(cfg)

	coinMarketCapProvider := coinMarketCapProviderCreator.CreateCryptoProvider()
	binanceProvider := binanceProviderCreator.CreateCryptoProvider()
	coinAPIProvider := coinAPIProviderCreator.CreateCryptoProvider()
	coinbaseProvider := coinbaseProviderCreator.CreateCryptoProvider()

	coinMarketCapChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(coinMarketCapProvider, logger))
	binanceChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(binanceProvider, logger))
	coinAPIChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(coinAPIProvider, logger))
	coinbaseChain := details.NewBaseCryptoChain(crypto_providers.NewLoggingCryptoProvider(coinbaseProvider, logger))

	coinMarketCapChain.SetNext(binanceChain)
	binanceChain.SetNext(coinAPIChain)
	coinAPIChain.SetNext(coinbaseChain)

	return coinMarketCapChain
}

func initRepos(filePath string, logger *logrus.Logger) *usecases_contracts.Repository {
	emailSub := filestorage.NewEmailSubscriptionRepository(filePath, logger)

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
