package main

import (
	"os"

	"user-auth/middleware"
	routes "user-auth/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	router.Static("/assets", "./assets")
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.NoteRoutes(router)

	router.Run(":" + port)
}
