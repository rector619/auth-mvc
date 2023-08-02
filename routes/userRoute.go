package routes

import (
	controller "github.com/aremxyplug-be/controllers"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.POST("/api/v1/users/signup", controller.SignUp())
    incomingRoutes.POST("/api/v1/users/login", controller.Login())
}