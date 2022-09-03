package handler

import (
	"api/internal/service"
	mock_service "api/internal/service/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHTTPHandler_getCurrentExchangeRate(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCrypto)

	type test struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}

	testTable := []test{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockCrypto) {
				s.EXPECT().GetBtcUahRate().Return(777.777, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "777.777",
		},
		{
			name: "Error",
			mockBehavior: func(s *mock_service.MockCrypto) {
				s.EXPECT().GetBtcUahRate().Return(float64(0), errors.New("some error"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)

			cryptoServiceMock := mock_service.NewMockCrypto(mockController)
			testCase.mockBehavior(cryptoServiceMock)

			services := &service.Service{Crypto: cryptoServiceMock}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/rate", handler.getCurrentExchangeRate)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/rate", nil)

			r.ServeHTTP(responseRecorder, request)

			assert.Equal(t, testCase.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, responseRecorder.Body.String())

		})
	}

}