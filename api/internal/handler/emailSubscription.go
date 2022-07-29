package handler

import (
	"api/internal/command"
	"api/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) sendEmails(c *gin.Context) {
	err := h.services.EmailSub.SendToAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "sent"})
}

func (h *Handler) subscribe(c *gin.Context) {
	var input command.Subscribing

	err := c.ShouldBind(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Email field is required")
		return
	}

	err = h.services.Subscribe(input.Email)
	if err != nil {
		if errors.Is(err, service.EmailDuplErr) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "subscribed"})
}
