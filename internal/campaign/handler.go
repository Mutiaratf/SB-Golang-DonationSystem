package campaign

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

// @Summary Create campaign
// @Tags Campaigns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CampaignRequest true "Campaign"
// @Success 201 {object} map[string]interface{}
// @Failure 400,401,500 {object} map[string]string
// @Router /campaigns [post]
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

// @Summary List campaigns
// @Tags Campaigns
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /campaigns [get]
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

// @Summary Get campaign
// @Tags Campaigns
// @Produce json
// @Param id path int true "Campaign ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]string
// @Router /campaigns/{id} [get]
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

// @Summary Update campaign
// @Tags Campaigns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Campaign ID"
// @Param request body CampaignRequest true "Campaign"
// @Success 200 {object} map[string]interface{}
// @Failure 400,401,404,500 {object} map[string]string
// @Router /campaigns/{id} [put]
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

// @Summary Delete campaign
// @Tags Campaigns
// @Security BearerAuth
// @Produce json
// @Param id path int true "Campaign ID"
// @Success 200 {object} map[string]string
// @Failure 400,401,404,500 {object} map[string]string
// @Router /campaigns/{id} [delete]
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
	if errors.Is(err, ErrCampaignNameEmpty) || errors.Is(err, ErrCategoryIDInvalid) || errors.Is(err, ErrAmountInvalid) || errors.Is(err, ErrCategoryNotFound) || errors.Is(err, ErrCategoryInactive) {
		status = http.StatusBadRequest
		message = err.Error()
	} else if errors.Is(err, ErrCampaignNotFound) {
		status = http.StatusNotFound
		message = "Campaign not found"
	}
	c.JSON(status, gin.H{"status": "error", "message": message})
}
