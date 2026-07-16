package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/database"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

//go:embed templates/*.html
var fsys embed.FS

func main() {
	router := gin.Default()

	templ := template.Must(template.New("").ParseFS(fsys, "templates/*.html"))
	router.SetHTMLTemplate(templ)
	router.Static("/static", "./static")

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	var client *mongo.Client = database.Connect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to reach server: %v", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "BrickBase - Home",
		})
	})

	router.GET("/sets", func(c *gin.Context) {
		c.HTML(http.StatusOK, "set.html", gin.H{
			"title": "Catalogo Set LEGO",
		})
	})

	routes.SetupUnProtectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
