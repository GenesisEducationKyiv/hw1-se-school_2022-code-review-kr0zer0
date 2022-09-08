package handler

import (
	"api/internal/customerrors"
	"api/internal/inputs"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) sendEmails(c *gin.Context) {
	err := h.services.EmailSubService.SendToAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "sent"})
}

func (h *HTTPHandler) subscribe(c *gin.Context) {
	var input inputs.Subscribing

	err := c.ShouldBind(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Email field is required")
		return
	}

	err = h.services.Subscribe(input.Email)
	if err != nil {
		if errors.Is(err, customerrors.ErrEmailDuplicate) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "subscribed"})
}
