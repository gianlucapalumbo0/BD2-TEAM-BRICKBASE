package main

import (
	"fmt"

	controller "github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	router.GET("/sets", controller.GetSets())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
