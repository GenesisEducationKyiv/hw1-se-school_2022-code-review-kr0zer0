package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getCurrentExchangeRate(c *gin.Context) {
	rate, err := h.services.Crypto.GetCurrentExchangeRate()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, rate)
}
