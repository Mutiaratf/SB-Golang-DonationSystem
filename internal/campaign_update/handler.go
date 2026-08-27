package campaign_update

import (
	"errors"
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

// @Summary List campaign updates
// @Tags Campaign Updates
// @Produce json
// @Param id path int true "Campaign ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]string
// @Router /campaigns/{id}/updates [get]
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

// @Summary Create campaign update
// @Tags Campaign Updates
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Campaign ID"
// @Param request body CampaignUpdateRequest true "Campaign update"
// @Success 201 {object} map[string]interface{}
// @Failure 400,401,500 {object} map[string]string
// @Router /campaigns/{id}/updates [post]
func (h *Handler) Create(c *gin.Context) {
	campaignID, ok := parseID(c, "id")
	if !ok {
		return
	}
	var request CampaignUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "title dan content wajib diisi")
		return
	}
	update, err := h.service.Create(campaignID, &request)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Campaign update created successfully", "data": update})
}

// @Summary Update campaign update
// @Tags Campaign Updates
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Campaign ID"
// @Param update_id path int true "Campaign update ID"
// @Param request body CampaignUpdateRequest true "Campaign update"
// @Success 200 {object} map[string]interface{}
// @Failure 400,401,404,500 {object} map[string]string
// @Router /campaigns/{id}/updates/{update_id} [put]
func (h *Handler) Update(c *gin.Context) {
	campaignID, ok := parseID(c, "id")
	if !ok {
		return
	}
	updateID, ok := parseID(c, "update_id")
	if !ok {
		return
	}
	var request CampaignUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "title dan content wajib diisi")
		return
	}
	update, err := h.service.Update(updateID, campaignID, &request)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Campaign update updated successfully", "data": update})
}

// @Summary Delete campaign update
// @Tags Campaign Updates
// @Security BearerAuth
// @Produce json
// @Param id path int true "Campaign ID"
// @Param update_id path int true "Campaign update ID"
// @Success 200 {object} map[string]string
// @Failure 400,401,404,500 {object} map[string]string
// @Router /campaigns/{id}/updates/{update_id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	campaignID, ok := parseID(c, "id")
	if !ok {
		return
	}
	updateID, ok := parseID(c, "update_id")
	if !ok {
		return
	}
	if err := h.service.Delete(updateID, campaignID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Campaign update deleted successfully"})
}

func parseID(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "Invalid "+name)
		return 0, false
	}
	return id, true
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrCampaignNotFound):
		respondError(c, http.StatusBadRequest, "campaign tidak ditemukan")
	case errors.Is(err, ErrCampaignInactive):
		respondError(c, http.StatusBadRequest, "campaign tidak aktif")
	case errors.Is(err, ErrCampaignUpdateNotFound):
		respondError(c, http.StatusNotFound, "campaign update tidak ditemukan")
	default:
		respondError(c, http.StatusInternalServerError, "Failed to process campaign update")
	}
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"status": "error", "message": message})
}
