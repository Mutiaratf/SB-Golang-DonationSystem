package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service *Service
	config  config.Config
}

func NewHandler(service *Service, cfg config.Config) *Handler {
	return &Handler{service: service, config: cfg}
}

// Login authenticates a user and returns a JWT.
// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400,401,500 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid login request"})
		return
	}

	account, err := h.service.Authenticate(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid email or password"})
		return
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    h.config.JWTIssuer,
		Subject:   strconv.FormatInt(account.ID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(h.config.JWTExpiry)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   LoginResponse{Token: signedToken, User: account},
	})
}

// Register creates a user account.
// @Summary Register
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{}
// @Failure 400,409,500 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email harus valid dan password minimal 8 karakter"})
		return
	}

	account, err := h.service.Register(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Email sudah terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat akun"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registration successful",
		"data":    account,
	})
}
