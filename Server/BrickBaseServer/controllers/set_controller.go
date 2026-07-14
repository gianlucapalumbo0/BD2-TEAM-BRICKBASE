package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/database"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/models"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/api/option"
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
			return
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
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can add sets"})
			return
		}

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
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can add sets"})
			return
		}

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
				"name":            set.Name,
				"year":            set.Year,
				"theme_id":        set.ThemeID,
				"num_parts":       set.NumParts,
				"parts_inventory": set.PartsInventory,
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

		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can add sets"})
			return
		}

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

func AddUserReview() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}

		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}

		if role != "USER" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Only users can add reviews",
			})
			return
		}

		setId := c.Param("set_num")
		if setId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Set Id required"})
			return
		}

		var req struct {
			Review string `json:"review"`
		}

		var resp struct {
			UserID string  `json:"user_id"`
			Review string  `json:"review"`
			Rating float64 `json:"rating"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setCollection := database.OpenCollection("sets", database.Connect())

		checkFilter := bson.M{
			"set_num":              setId,
			"user_reviews.user_id": userId,
		}

		var existingSet models.Set

		err = setCollection.FindOne(ctx, checkFilter).Decode(&existingSet)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already reviewed this set",
			})
			return
		}

		rating, err := GetReviewRating(req.Review, database.Connect(), c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filter := bson.D{
			{Key: "set_num", Value: setId},
		}

		update := bson.M{
			"$push": bson.M{
				"user_reviews": models.UserReview{
					UserID: userId,
					Review: req.Review,
					Rating: rating,
				},
			},
		}

		result, err := setCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding review"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Set not found"})
			return
		}

		var set models.Set

		err = setCollection.FindOne(ctx, filter).Decode(&set)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error retrieving set",
			})
			return
		}

		var total float64

		for _, review := range set.UserReviews {
			total += review.Rating
		}

		average := total / float64(len(set.UserReviews))

		_, err = setCollection.UpdateOne(
			ctx,
			filter,
			bson.M{
				"$set": bson.M{
					"review_rating": average,
				},
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error updating review rating",
			})
			return
		}

		resp.UserID = userId
		resp.Review = req.Review
		resp.Rating = rating

		c.JSON(http.StatusOK, resp)
	}
}

func GetReviewRating(review string, client *mongo.Client, c *gin.Context) (float64, error) {

	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	geminiApiKey := os.Getenv("GEMINI_API_KEY")

	if geminiApiKey == "" {
		return 0, errors.New("could not read GEMINI_API_KEY")
	}

	ctx := context.Background()

	geminiClient, err := genai.NewClient(
		ctx,
		option.WithAPIKey(geminiApiKey),
	)

	if err != nil {
		return 0, err
	}

	defer geminiClient.Close()

	model := geminiClient.GenerativeModel("gemini-flash-latest")

	basePrompt := os.Getenv("BASE_PROMPT_TEMPLATE")

	response, err := model.GenerateContent(
		ctx,
		genai.Text(basePrompt+review),
	)

	if err != nil {
		return 0, err
	}

	textResponse := response.Candidates[0].Content.Parts[0].(genai.Text)

	rating, err := strconv.ParseFloat(string(textResponse), 64)

	if err != nil {
		return 0, errors.New("invalid rating returned by Gemini")
	}

	if rating < 1 || rating > 5 {
		return 0, errors.New("rating must be between 1 and 5")
	}

	return rating, nil
}

func GetBestSets() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env file not found")
		}

		var limit int64 = 5

		findOptions := options.Find()

		findOptions.SetSort(
			bson.D{
				{
					Key:   "review_rating",
					Value: -1,
				},
			},
		)

		findOptions.SetLimit(limit)

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setCollection := database.OpenCollection("sets", database.Connect())

		cursor, err := setCollection.Find(
			ctx,
			bson.D{},
			findOptions,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching best sets",
			})
			return
		}

		defer cursor.Close(ctx)

		var bestSets []models.Set

		if err := cursor.All(ctx, &bestSets); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, bestSets)
	}
}
