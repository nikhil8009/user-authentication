package routes

import (
	controller "user-auth/controllers"
	"user-auth/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/user/:user_id", controller.GetUser())
	incomingRoutes.PUT("/user/:user_id", controller.UpdateUser())
	incomingRoutes.DELETE("/user/:user_id", controller.DeleteUser())
	incomingRoutes.POST("/users/upload", controller.UploadFile())
	incomingRoutes.POST("/users/updateImage", controller.UpdateImage())
}
