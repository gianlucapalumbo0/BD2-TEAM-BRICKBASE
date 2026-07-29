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
func GetSets(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		skip := (page - 1) * limit

		var setCollection *mongo.Collection = database.OpenCollection("sets", client)
		var sets []models.Set

		findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))

		cursor, err := setCollection.Find(ctx, bson.M{}, findOptions)
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
func GetSet(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setID := c.Param("set_num")
		fmt.Println("SET ID:", setID)

		if setID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Set ID is required"})
			return
		}

		var setCollection *mongo.Collection = database.OpenCollection("sets", client)

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
func AddSet(client *mongo.Client) gin.HandlerFunc {
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
		var setCollection *mongo.Collection = database.OpenCollection("sets", client)

		result, err := setCollection.InsertOne(ctx, set)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add set"})
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}

// UpdateSet aggiorna un set esistente dato il suo numero
func UpdateSet(client *mongo.Client) gin.HandlerFunc {
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

		var setCollection *mongo.Collection = database.OpenCollection("sets", client)

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
func DeleteSet(client *mongo.Client) gin.HandlerFunc {
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

		var setCollection *mongo.Collection = database.OpenCollection("sets", client)

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

func AddUserReview(client *mongo.Client) gin.HandlerFunc {
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
			UserID       string  `json:"user_id"`
			Review       string  `json:"review"`
			Rating       float64 `json:"rating"`
			ReviewRating float64 `json:"review_rating"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setCollection := database.OpenCollection("sets", client)

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

		rating, err := GetReviewRating(req.Review, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filter := bson.D{
			{Key: "set_num", Value: setId},
		}

		// CORREZIONE 1: Rimosso "ReviewRating: average" da qui (average non esiste ancora)
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

		// Qui viene calcolata la media
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

		// CORREZIONE 2: Assegnazione della media calcolata al JSON di risposta
		resp.UserID = userId
		resp.Review = req.Review
		resp.Rating = rating
		resp.ReviewRating = average

		c.JSON(http.StatusOK, resp)
	}
}

func GetReviewRating(review string, client *mongo.Client, c *gin.Context) (float64, error) {

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

	if rating < 1 {
		rating = 1
	} else if rating > 5 {
		rating = 5
	}

	return rating, nil
}

func GetBestSets(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env file not found")
		}

		var limit int64 = 4

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

		setCollection := database.OpenCollection("sets", client)

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

// GetSetInventoryDetailed restituisce tutti i pezzi di un set completi di immagini e colori
func GetSetInventoryDetailed(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setID := c.Param("set_num")
		if setID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Set ID is required"})
			return
		}

		var setCollection *mongo.Collection = database.OpenCollection("sets", client)

		// Pipeline di Aggregazione MongoDB
		pipeline := mongo.Pipeline{
			// 1. Trova il set specifico
			{{Key: "$match", Value: bson.M{"set_num": setID}}},

			// 2. "Srotola" l'array dell'inventario in documenti separati
			{{Key: "$unwind", Value: "$parts_inventory"}},

			// 3. Unisci i dati del pezzo (collezione parts)
			{{Key: "$lookup", Value: bson.M{
				"from":         "parts",
				"localField":   "parts_inventory.part_num",
				"foreignField": "part_num",
				"as":           "part_info",
			}}},
			// Trasforma l'array risultante in un oggetto
			{{Key: "$unwind", Value: bson.M{"path": "$part_info", "preserveNullAndEmptyArrays": true}}},

			// 4. Unisci i dati del colore (collezione colors)
			{{Key: "$lookup", Value: bson.M{
				"from":         "colors",
				"localField":   "parts_inventory.color_id",
				"foreignField": "color_id",
				"as":           "color_info",
			}}},
			{{Key: "$unwind", Value: bson.M{"path": "$color_info", "preserveNullAndEmptyArrays": true}}},

			// 5. Proietta (seleziona) solo i campi che ci interessano per il frontend
			{{Key: "$project", Value: bson.M{
				"_id":          0,
				"part_num":     "$parts_inventory.part_num",
				"quantity":     "$parts_inventory.quantity",
				"is_spare":     "$parts_inventory.is_spare",
				"part_name":    "$part_info.name",
				"part_img_url": "$part_info.part_img_url",
				"color_name":   "$color_info.name",
				"color_rgb":    "$color_info.rgb",
			}}},
		}

		cursor, err := setCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero dei pezzi"})
			return
		}
		defer cursor.Close(ctx)

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore decodifica pezzi"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// GetUserReviewedSets restituisce un array di set_num che l'utente loggato ha già recensito
func GetUserReviewedSets(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recuperiamo l'ID dell'utente dal token/contesto
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Utente non autorizzato"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		setCollection := database.OpenCollection("sets", client)

		// Filtriamo i set in cui l'array user_reviews contiene un oggetto col nostro user_id
		filter := bson.M{"user_reviews.user_id": userId}

		// Proiezione: ci interessa SOLO il set_num (per risparmiare banda)
		findOptions := options.Find().SetProjection(bson.M{"set_num": 1, "_id": 0})

		cursor, err := setCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore nel recupero delle recensioni"})
			return
		}
		defer cursor.Close(ctx)

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore decodifica risultati"})
			return
		}

		// Estraiamo solo le stringhe set_num
		reviewedSetNums := make([]string, 0)
		for _, res := range results {
			if setNum, ok := res["set_num"].(string); ok {
				reviewedSetNums = append(reviewedSetNums, setNum)
			}
		}

		c.JSON(http.StatusOK, reviewedSetNums)
	}
}
