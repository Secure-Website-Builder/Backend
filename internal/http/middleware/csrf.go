package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var csrfMethods = map[string]bool{
	http.MethodPost:   true,
	http.MethodPut:    true,
	http.MethodPatch:  true,
	http.MethodDelete: true,
}

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		if !csrfMethods[c.Request.Method] {
			c.Next()
			return
		}

		if auth := c.GetHeader("Authorization"); auth != "" {
			c.Next()
			return
		}

		headerToken := c.GetHeader("X-CSRF-Token")
		cookieToken, err := c.Cookie("csrf_token")

		if err != nil || headerToken == "" || headerToken != cookieToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token mismatch or missing",
			})
			return
		}

		c.Next()
	}
}


func GenerateCSRFToken() string {
	return uuid.NewString()
}
