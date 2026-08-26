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
)

// SearchThemes cerca i temi nel database per nome (case-insensitive)
func SearchThemes(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusOK, []models.Theme{})
			return
		}

		ctx, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()

		var themeCollection *mongo.Collection = database.OpenCollection("themes", client)

		// Ricerca testuale con regex case-insensitive ("$options": "i")
		filter := bson.M{"name": bson.M{"$regex": query, "$options": "i"}}

		cursor, err := themeCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante la ricerca dei temi"})
			return
		}
		defer cursor.Close(ctx)

		var themes []models.Theme
		if err = cursor.All(ctx, &themes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore nel parsing dei temi"})
			return
		}

		if themes == nil {
			themes = []models.Theme{}
		}

		c.JSON(http.StatusOK, themes)
	}
}
