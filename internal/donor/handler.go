package donor

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	donorService *Service
}

func NewHandler(donorService *Service) *Handler {
	return &Handler{donorService: donorService}
}

func (h *Handler) CreateDonor(c *gin.Context) {
	var req DonorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}
	donor, err := h.donorService.CreateDonor(&req)
	if err != nil {
		handleDonorError(c, err, "Failed to create donor")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Donor created successfully",
		"data":    donor,
	})
}

func (h *Handler) GetAllDonors(c *gin.Context) {
	donors, err := h.donorService.GetAllDonors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve donors",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Donors retrieved successfully",
		"data":    donors,
	})
}

func (h *Handler) GetDonorByID(c *gin.Context) {
	id, ok := donorID(c)
	if !ok {
		return
	}
	donor, err := h.donorService.GetDonorByID(id)
	if err != nil {
		handleDonorError(c, err, "Failed to retrieve donor")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Donor retrieved successfully",
		"data":    donor,
	})
}

func (h *Handler) UpdateDonor(c *gin.Context) {
	id, ok := donorID(c)
	if !ok {
		return
	}
	var req DonorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}
	donor, err := h.donorService.UpdateDonor(id, &req)
	if err != nil {
		handleDonorError(c, err, "Failed to update donor")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Donor updated successfully",
		"data":    donor,
	})
}

func (h *Handler) DeleteDonor(c *gin.Context) {
	id, ok := donorID(c)
	if !ok {
		return
	}
	if err := h.donorService.DeleteDonor(id); err != nil {
		handleDonorError(c, err, "Failed to delete donor")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Donor deleted successfully"})
}

func donorID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid donor ID",
		})
		return 0, false
	}
	return id, true
}

func handleDonorError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, ErrDonorNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Donor not found",
		})
		return
	}
	if errors.Is(err, ErrDonorNameEmpty) || errors.Is(err, ErrDonorEmailEmpty) ||
		errors.Is(err, ErrDonorEmailInvalid) || errors.Is(err, ErrDonorPhoneEmpty) ||
		errors.Is(err, ErrDonorPhoneInvalid) || errors.Is(err, ErrDonorGenderInvalid) ||
		errors.Is(err, ErrDonorPhoneExists) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": fallback,
	})
}
