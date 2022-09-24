package integration

import (
	"api/config"
	"api/internal/constants"
	"api/internal/controllers/http"
	"api/internal/infrastructure/cryptoProviders"
	"api/internal/infrastructure/repository/file"
	"api/internal/service"
	"api/internal/service/crypto"
	"api/internal/service/interfaces"
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
	cryptoChain crypto.CryptoChain
	handler     *http.Handler
	services    *interfaces.Service
	repos       *interfaces.Repository

	mailerMock *mock_service.MockMailer
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
	s.mailerMock = mock_service.NewMockMailer(mockController)

	s.cfg = config.GetConfig()
	coinMarketCapProviderCreator := crypto_providers.NewCoinMarketCapProviderCreator(s.cfg)
	binanceProviderCreator := crypto_providers.NewBinanceProviderCreator(s.cfg)
	coinAPIProviderCreator := crypto_providers.NewCoinAPIProviderCreator(s.cfg)
	coinbaseProviderCreator := crypto_providers.NewCoinbaseProviderCreator(s.cfg)

	coinMarketCapProvider := coinMarketCapProviderCreator.CreateCryptoProvider()
	binanceProvider := binanceProviderCreator.CreateCryptoProvider()
	coinAPIProvider := coinAPIProviderCreator.CreateCryptoProvider()
	coinbaseProvider := coinbaseProviderCreator.CreateCryptoProvider()

	coinMarketCapChain := crypto.NewBaseCryptoChain(coinMarketCapProvider)
	binanceChain := crypto.NewBaseCryptoChain(binanceProvider)
	coinAPIChain := crypto.NewBaseCryptoChain(coinAPIProvider)
	coinbaseChain := crypto.NewBaseCryptoChain(coinbaseProvider)

	coinMarketCapChain.SetNext(binanceChain)
	binanceChain.SetNext(coinAPIChain)
	coinAPIChain.SetNext(coinbaseChain)

	s.cryptoChain = coinMarketCapChain
	s.repos = &interfaces.Repository{
		EmailSubscriptionRepo: file.NewEmailSubscriptionRepository(TestDataPath),
	}
	s.services = service.NewService(s.repos, s.cryptoChain, s.mailerMock, s.cfg)
	s.handler = http.NewHandler(s.services)
}
