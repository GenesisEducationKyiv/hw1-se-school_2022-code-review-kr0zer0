package integration

import (
	"api/config"
	"api/internal/constants"
	"api/internal/handler"
	"api/internal/repository"
	"api/internal/service"
	"api/internal/service/cryptoProviders"
	mock_service "api/internal/service/mocks"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

const TestDataPath = "../data/data.json"

type IntegrationTestSuite struct {
	suite.Suite

	cfg         *config.Config
	cryptoChain service.CryptoChain
	handler     *handler.HTTPHandler
	services    *handler.Service
	repos       *service.Repository

	emailSendingMock *mock_service.MockEmailSendingRepo
}

func TestIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		panic(err)
	}

	dataToWrite := []byte(`{"emails":[]}`)
	err = os.WriteFile(TestDataPath, dataToWrite, constants.WriteFilePerm)
	if err != nil {
		s.FailNowf("unable to setup data file", err.Error())
	}

	s.initDeps()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	err := os.Truncate(TestDataPath, 0)
	if err != nil {
		s.FailNowf("unable to clear data file", err.Error())
	}
}

func (s *IntegrationTestSuite) initDeps() {
	mockController := gomock.NewController(s.T())
	s.emailSendingMock = mock_service.NewMockEmailSendingRepo(mockController)

	s.cfg = config.GetConfig()
	coinMarketCapProviderCreator := crypto_providers.NewCoinMarketCapProviderCreator(s.cfg)
	binanceProviderCreator := crypto_providers.NewBinanceProviderCreator(s.cfg)
	coinAPIProviderCreator := crypto_providers.NewCoinAPIProviderCreator(s.cfg)
	coinbaseProviderCreator := crypto_providers.NewCoinbaseProviderCreator(s.cfg)

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

	s.cryptoChain = coinMarketCapChain
	s.repos = &service.Repository{
		EmailSubscriptionRepo: repository.NewEmailSubscriptionRepository(TestDataPath),
		EmailSendingRepo:      s.emailSendingMock,
	}
	s.services = service.NewService(s.repos, s.cryptoChain, s.cfg)
	s.handler = handler.NewHandler(s.services)
}
