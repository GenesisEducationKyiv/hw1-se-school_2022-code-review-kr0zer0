package integration

import (
	"api/config"
	"api/internal/constants"
	"api/internal/handler"
	"api/internal/repository"
	mock_repository "api/internal/repository/mocks"
	"api/internal/service"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

const TestDataPath = "../data/data.json"

type IntegrationTestSuite struct {
	suite.Suite

	cfg      *config.Config
	handler  *handler.HTTPHandler
	services *handler.Service
	repos    *service.Repository

	emailSendingMock *mock_repository.MockEmailSending
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
	s.emailSendingMock = mock_repository.NewMockEmailSending(mockController)

	s.cfg = config.GetConfig()
	s.repos = &service.Repository{
		EmailSubscriptionRepo: repository.NewEmailSubscriptionRepository(TestDataPath),
		EmailSendingRepo:      s.emailSendingMock,
	}
	s.services = service.NewService(s.repos, s.cfg)
	s.handler = handler.NewHandler(s.services)
}
