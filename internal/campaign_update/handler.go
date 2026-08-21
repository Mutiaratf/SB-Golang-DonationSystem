package campaign_update

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByCampaignID(c *gin.Context) {
	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil || campaignID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid campaign ID",
		})
		return
	}

	updates, err := h.service.GetByCampaignID(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve campaign updates",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Campaign updates retrieved successfully",
		"data":    updates,
	})
}
