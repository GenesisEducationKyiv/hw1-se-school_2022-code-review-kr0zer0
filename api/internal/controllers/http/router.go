package http

import (
	"api/internal/service"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=router.go -destination=mocks/serviceMock.go

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()

	base := router.Group("/api")
	base.GET("/rate", h.getCurrentExchangeRate)
	base.POST("/subscribe", h.subscribe)
	base.POST("/sendEmails", h.sendEmails)

	return router
}
