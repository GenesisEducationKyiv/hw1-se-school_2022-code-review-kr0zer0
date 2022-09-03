package handler

import (
	"api/internal/customerrors"
	"api/internal/service"
	mock_service "api/internal/service/mocks"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHTTPHandler_sendMails(t *testing.T) {
	type mockBehavior func(s *mock_service.MockEmailSub)

	type test struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}

	testTable := []test{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockEmailSub) {
				s.EXPECT().SendToAll().Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"sent"}`,
		},
		{
			name: "Error",
			mockBehavior: func(s *mock_service.MockEmailSub) {
				s.EXPECT().SendToAll().Return(errors.New("some error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)

			emailSubMock := mock_service.NewMockEmailSub(mockController)
			testCase.mockBehavior(emailSubMock)

			services := &service.Service{EmailSub: emailSubMock}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/send", handler.sendEmails)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/send", bytes.NewBufferString(""))

			r.ServeHTTP(responseRecorder, request)

			assert.Equal(t, testCase.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, responseRecorder.Body.String())

		})
	}
}

func TestHTTPHandler_subscribe(t *testing.T) {
	type mockBehavior func(s *mock_service.MockEmailSub, email string)

	type test struct {
		name                 string
		emailInput           string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}

	testTable := []test{
		{
			name:       "OK",
			emailInput: "some.email@mail.com",
			mockBehavior: func(s *mock_service.MockEmailSub, email string) {
				s.EXPECT().Subscribe(email).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"subscribed"}`,
		},
		{
			name: "No email input",
			mockBehavior: func(s *mock_service.MockEmailSub, email string) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Email field is required"}`,
		},
		{
			name:       "Email duplicate",
			emailInput: "some.email@mail.com",
			mockBehavior: func(s *mock_service.MockEmailSub, email string) {
				s.EXPECT().Subscribe(email).Return(customerrors.ErrEmailDuplicate)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: fmt.Sprintf(`{"message":"%s"}`, customerrors.ErrEmailDuplicate.Error()),
		},
		{
			name:       "Some internal error",
			emailInput: "some.email@mail.com",
			mockBehavior: func(s *mock_service.MockEmailSub, email string) {
				s.EXPECT().Subscribe(email).Return(errors.New("some error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)

			emailSubMock := mock_service.NewMockEmailSub(mockController)
			testCase.mockBehavior(emailSubMock, testCase.emailInput)

			services := &service.Service{EmailSub: emailSubMock}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/subscribe", handler.subscribe)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/subscribe?email="+testCase.emailInput, nil)

			r.ServeHTTP(responseRecorder, request)

			assert.Equal(t, testCase.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, responseRecorder.Body.String())

		})
	}
}