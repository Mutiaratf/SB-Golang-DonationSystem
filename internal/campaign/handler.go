package campaign

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreateCampaign(c *gin.Context) {
	var request CampaignRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}
	campaign, err := h.service.Create(&request)
	if err != nil {
		h.respondError(c, err, "create")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Campaign created successfully", "data": campaign})
}

func (h *Handler) GetAllCampaigns(c *gin.Context) {
	campaigns, err := h.service.GetAll()
	if err != nil {
		h.respondError(c, err, "retrieve")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Campaigns retrieved successfully",
		"data":    campaigns,
	})
}

func (h *Handler) GetCampaignByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	campaign, err := h.service.GetByID(id)
	if err != nil {
		h.respondError(c, err, "retrieve")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Campaign retrieved successfully",
		"data":    campaign,
	})
}

func (h *Handler) UpdateCampaign(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var request CampaignRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}
	campaign, err := h.service.Update(id, &request)
	if err != nil {
		h.respondError(c, err, "update")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Campaign updated successfully",
		"data":    campaign,
	})
}

func (h *Handler) DeleteCampaign(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		h.respondError(c, err, "delete")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Campaign deleted successfully",
	})
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid campaign ID"})
		return 0, false
	}
	return id, true
}

func (h *Handler) respondError(c *gin.Context, err error, action string) {
	status := http.StatusInternalServerError
	message := "Failed to " + action + " campaign"
	if errors.Is(err, ErrCampaignNameEmpty) || errors.Is(err, ErrCategoryIDInvalid) || errors.Is(err, ErrAmountInvalid) || errors.Is(err, ErrCategoryInactive) {
		status = http.StatusBadRequest
		message = err.Error()
	} else if errors.Is(err, ErrCampaignNotFound) {
		status = http.StatusNotFound
		message = "Campaign not found"
	}
	c.JSON(status, gin.H{"status": "error", "message": message})
}
