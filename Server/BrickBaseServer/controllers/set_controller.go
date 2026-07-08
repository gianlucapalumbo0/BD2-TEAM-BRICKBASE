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

var setCollection *mongo.Collection = database.OpenCollection("sets", database.Connect())

func GetSets() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var sets []models.Set

		cursor, err := setCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sets."})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &sets); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode sets."})
			return
		}

		c.JSON(http.StatusOK, sets)

	}
}
