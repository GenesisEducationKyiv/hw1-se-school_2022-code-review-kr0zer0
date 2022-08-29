package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getCurrentExchangeRate(c *gin.Context) {
	rate, err := h.services.Crypto.GetBtcUahRate()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, rate)
}
