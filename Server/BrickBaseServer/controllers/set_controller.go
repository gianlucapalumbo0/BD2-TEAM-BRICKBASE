package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/database"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

// restituisce tutti i set presenti nel database
func GetSets() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var setCollection *mongo.Collection = database.OpenCollection("sets", database.Connect())

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

// restituisce un set specifico dato il suo numero
func GetSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setID := c.Param("set_num")
		fmt.Println("SET ID:", setID)

		if setID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Set ID is required"})
			return
		}

		var setCollection *mongo.Collection = database.OpenCollection("sets", database.Connect())

		var set models.Set

		err := setCollection.FindOne(ctx, bson.M{"set_num": setID}).Decode(&set)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Set not found"})
			return
		}

		c.JSON(http.StatusOK, set)

	}
}

// AddSet inserisce un nuovo set nel database
func AddSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var set models.Set
		if err := c.ShouldBindJSON(&set); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(set); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		var setCollection *mongo.Collection = database.OpenCollection("sets", database.Connect())

		result, err := setCollection.InsertOne(ctx, set)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add set"})
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}

// UpdateSet aggiorna un set esistente dato il suo numero
func UpdateSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setID := c.Param("set_num")

		if setID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Set ID is required"})
			return
		}

		var set models.Set
		if err := c.ShouldBindJSON(&set); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(set); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		var setCollection *mongo.Collection = database.OpenCollection("sets", database.Connect())

		update := bson.M{
			"$set": bson.M{
				"name":      set.Name,
				"year":      set.Year,
				"theme_id":  set.ThemeID,
				"num_parts": set.NumParts,
			},
		}

		result, err := setCollection.UpdateOne(ctx, bson.M{"set_num": setID}, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update set"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Set not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Set updated successfully"})
	}
}

// DeleteSet elimina un set dato il suo numero
func DeleteSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setID := c.Param("set_num")

		if setID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Set ID is required"})
			return
		}

		var setCollection *mongo.Collection = database.OpenCollection("sets", database.Connect())

		result, err := setCollection.DeleteOne(ctx, bson.M{"set_num": setID})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete set"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Set not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Set deleted successfully"})
	}
}
