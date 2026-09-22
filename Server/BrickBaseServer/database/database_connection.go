package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {

	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: unable to fund .env file")
	}

	MongoDb := os.Getenv("MONGODB_URI")

	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set!")
	}

	clientOptions := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil
	}

	return client
}

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Fatal("DATABASE_NAME not set!")
	}

	collection := client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}
	return collection

}

// CreateIndexes viene chiamata una volta all'avvio per configurare gli indici
func CreateIndexes(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	setCollection := OpenCollection("sets", client)

	// Definisce l'indice su review_rating in ordine decrescente (-1)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "review_rating", Value: -1},
		},
	}

	// Crea l'indice su MongoDB
	name, err := setCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Errore durante la creazione dell'indice su review_rating: %v\n", err)
	} else {
		log.Printf("Indice verificato/creato con successo: %s\n", name)
	}
}
