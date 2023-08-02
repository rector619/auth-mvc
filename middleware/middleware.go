package middleware

import (
	"github.com/gin-gonic/gin"
)

// corsMiddleware handles the CORS middleware
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS",
		)
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Authorization, Origin, Accept, Content-Type, X-Requested-With, Access-Control-Request-Method, Access-Control-Request-Headers",
		)
		c.Writer.Header().Set("Access-control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return

		}
		c.Next()
	}
}

