package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthenticatedUserIDKey = "authenticated_user_id"
	AuthenticatedEmailKey  = "authenticated_email"
)

func JWT(secret, issuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		}, jwt.WithIssuer(issuer))
		if err != nil || !token.Valid {
			unauthorized(c)
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Subject == "" {
			unauthorized(c)
			return
		}
		c.Set(AuthenticatedUserIDKey, claims.Subject)
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authentication required"})
}
