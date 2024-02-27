package main

import (
	"os"

	routes "user-athentication/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Static("/assets", "./assets")
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.NoteRoutes(router)

	router.Run(":" + port)
}
