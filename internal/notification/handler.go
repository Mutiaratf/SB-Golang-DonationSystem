package notification

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) SendTransactionNotification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id_transaksi"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid transaction ID"})
		return
	}
	email, err := h.service.GetRecipientEmail(id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve transaction"})
		return
	}
	h.service.SendAsync(id)
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "Notification email is being sent",
		"email":   email,
	})
}
