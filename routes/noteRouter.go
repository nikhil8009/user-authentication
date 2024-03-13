package routes

import (
	controller "user-auth/controllers"
	"user-auth/middleware"

	"github.com/gin-gonic/gin"
)

func NoteRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.POST("/note", controller.AddNote())
	incomingRoutes.GET("/notes", controller.GetNotes())
	incomingRoutes.PUT("/note", controller.UpdateNote())
	incomingRoutes.DELETE("/note", controller.DeleteNote())
}
