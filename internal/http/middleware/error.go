package middleware

import (
	"github.com/Secure-Website-Builder/Backend/internal/errorx"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		httpErr := errorx.Resolve(err.Err)
		c.JSON(httpErr.Status, gin.H{"error": httpErr.Message})
	}
}
