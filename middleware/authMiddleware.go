package middleware

import (
	"fmt"
	"net/http"

	helper "github.com/aremxyplug-be/helpers"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientToken := c.Request.Header.Get("token")
        if clientToken == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
            c.Abort()
            return
        }

        claims, err := helper.ValidateToken(clientToken)
        if err != "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
            c.Abort()
            return
        }

        c.Set("email", claims.Email)
        c.Set("fullname", claims.Fullname)
        c.Set("username", claims.Username)
        c.Set("uid", claims.Uid)

        c.Next()

    }
}


