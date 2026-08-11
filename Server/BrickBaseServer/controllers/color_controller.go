package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/database"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SearchColors cerca i colori per nome nel database
func SearchColors(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusOK, []models.Color{})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		colorCollection := database.OpenCollection("colors", client)

		// Filtro per ricerca parziale (case-insensitive) sul nome del colore
		filter := bson.M{
			"name": bson.M{"$regex": query, "$options": "i"},
		}

		findOptions := options.Find().SetLimit(15)
		cursor, err := colorCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante la ricerca dei colori"})
			return
		}
		defer cursor.Close(ctx)

		var colors []models.Color
		if err = cursor.All(ctx, &colors); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore decodifica colori"})
			return
		}

		c.JSON(http.StatusOK, colors)
	}
}
