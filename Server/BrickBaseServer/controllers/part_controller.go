package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/database"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/models"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Struttura temporanea per leggere la risposta di Rebrickable
type RebrickablePartResponse struct {
	PartImgUrl string `json:"part_img_url"`
}

// SyncPartImages sincronizza le immagini mancanti da Rebrickable
func SyncPartImages(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Protezione: Solo admin dovrebbero poter lanciare questo script
		role, err := utils.GetRoleFromContext(c)
		if err != nil || role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Solo gli admin possono avviare la sincronizzazione"})
			return
		}

		apiKey := os.Getenv("REBRICKABLE_API_KEY") // Assicurati di averla nel .env
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API Key di Rebrickable mancante"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Sincronizzazione avviata in background. Ci vorrà del tempo!"})

		// Usiamo una Goroutine (background) per non bloccare la risposta del server
		// altrimenti la richiesta HTTP andrebbe in timeout dopo pochi secondi
		go func() {
			// Creiamo un nuovo context per la goroutine (senza timeout breve)
			ctx := context.Background()
			partCollection := database.OpenCollection("parts", client)

			// Ottimizzazione: cerchiamo solo i pezzi che NON hanno ancora un'immagine
			filter := bson.M{
				"$or": []bson.M{
					{"part_img_url": bson.M{"$exists": false}},
					{"part_img_url": ""},
				},
			}

			cursor, err := partCollection.Find(ctx, filter)
			if err != nil {
				log.Println("Errore ricerca parti senza immagine:", err)
				return
			}
			defer cursor.Close(ctx)

			var parts []models.Part
			if err = cursor.All(ctx, &parts); err != nil {
				log.Println("Errore decodifica parti:", err)
				return
			}

			totalParts := len(parts) // Salviamo il totale dei pezzi da elaborare
			log.Printf("Trovati %d pezzi da aggiornare. Inizio sincronizzazione...", totalParts)

			// Modifica: Aggiunto 'i' per tenere traccia dell'indice corrente nel ciclo
			for i, part := range parts {

				// Calcoliamo i pezzi rimanenti
				pezziRimanenti := totalParts - (i + 1)

				// Stampiamo il progresso prima di chiamare l'API
				log.Printf("[Progresso: %d/%d] Mancano %d pezzi - Controllo: %s", i+1, totalParts, pezziRimanenti, part.PartNum)

				// 1. Chiamiamo l'API
				url := fmt.Sprintf("https://rebrickable.com/api/v3/lego/parts/%s/?key=%s", part.PartNum, apiKey)
				resp, err := http.Get(url)

				if err != nil {
					log.Printf("Errore di rete per il pezzo %s", part.PartNum)
					time.Sleep(1200 * time.Millisecond)
					continue
				}

				// GESTIONE 404: Se il pezzo non esiste su Rebrickable
				if resp.StatusCode == 404 {
					log.Printf("Pezzo %s non trovato (404). Lo segno come NOT_FOUND.", part.PartNum)
					update := bson.M{"$set": bson.M{"part_img_url": "NOT_FOUND"}}
					partCollection.UpdateOne(ctx, bson.M{"_id": part.ID}, update)
					resp.Body.Close()
					time.Sleep(1200 * time.Millisecond)
					continue
				}

				// Gestione altri errori (es. 429, 500)
				if resp.StatusCode != 200 {
					log.Printf("Errore API per il pezzo %s (Status: %d)", part.PartNum, resp.StatusCode)
					resp.Body.Close()
					time.Sleep(1200 * time.Millisecond)
					continue
				}

				// 2. Decodifichiamo il JSON
				var apiResult RebrickablePartResponse
				if err := json.NewDecoder(resp.Body).Decode(&apiResult); err == nil && apiResult.PartImgUrl != "" {

					// 3. Salviamo l'URL nel database
					update := bson.M{"$set": bson.M{"part_img_url": apiResult.PartImgUrl}}
					_, err := partCollection.UpdateOne(ctx, bson.M{"_id": part.ID}, update)
					if err != nil {
						log.Printf("Errore salvataggio immagine per pezzo %s", part.PartNum)
					} else {
						log.Printf("Aggiornato pezzo %s", part.PartNum)
					}
				}
				resp.Body.Close()

				// IL PASSAGGIO CHIAVE: Aspettiamo 1.2 secondi
				time.Sleep(1200 * time.Millisecond)
			}

			log.Println("Sincronizzazione immagini completata!")
		}()
	}
}
