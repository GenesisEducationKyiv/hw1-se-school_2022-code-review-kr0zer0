package handler

import (
	"github.com/gin-gonic/gin"
)

type CryptoService interface {
	GetCurrentExchangeRate(cryptoSymbol, fiatSymbol string) (float64, error)
	GetBtcUahRate() (float64, error)
}

type EmailSubService interface {
	SendToAll() error
	Subscribe(email string) error
}

type Service struct {
	CryptoService
	EmailSubService
}

type HTTPHandler struct {
	services *Service
}

func NewHandler(services *Service) *HTTPHandler {
	return &HTTPHandler{
		services: services,
	}
}

func (h *HTTPHandler) InitRouter() *gin.Engine {
	router := gin.Default()

	base := router.Group("/api")
	base.GET("/rate", h.getCurrentExchangeRate)
	base.POST("/subscribe", h.subscribe)
	base.POST("/sendEmails", h.sendEmails)

	return router
}
