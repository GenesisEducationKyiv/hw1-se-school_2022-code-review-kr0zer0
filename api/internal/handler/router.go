package handler

import (
	"api/internal/service"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *HTTPHandler {
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

	//return router.Run(port)
	return router
}
