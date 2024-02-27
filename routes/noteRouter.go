package routes

import (
	controller "user-athentication/controllers"
	"user-athentication/middleware"

	"github.com/gin-gonic/gin"
)

func NoteRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.POST("/note", controller.AddNote())
	incomingRoutes.GET("/notes", controller.GetNotes())
}
