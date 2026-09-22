package models

import "go.mongodb.org/mongo-driver/v2/bson"

// InventoryPart rappresenta un elemento dell'inventario dei pezzi all'interno di un set
// PartNum è il codice alfanumerico del pezzo
// ColorID è l'identificativo numerico del colore del pezzo
// Quantity indica la quantità di questo specifico pezzo (con questo colore)
// IsSpare indica se il pezzo è fornito come pezzo di ricambio extra
type InventoryPart struct {
	PartNum  string `bson:"part_num" json:"part_num"`
	ColorID  int    `bson:"color_id" json:"color_id"`
	Quantity int    `bson:"quantity" json:"quantity"`
	IsSpare  bool   `bson:"is_spare" json:"is_spare"`
}

// UserReview rappresenta una singola recensione lasciata da un utente per un determinato Set
// UserID è l'identificativo univoco dell'utente che ha pubblicato la recensione
// Review contiene il testo della recensione
// Rating è la valutazione numerica assegnata alla recensione
type UserReview struct {
	UserID string  `bson:"user_id" json:"user_id"`
	Review string  `bson:"review" json:"review"`
	Rating float64 `bson:"rating" json:"rating"`
}

// Set rappresenta il documento principale all'interno della collezione "sets" in MongoDB
// ID è l'identificativo primario univoco generato da MongoDB
// SetNum è il codice univoco del set
// Name è il titolo del set
// Year è l'anno di uscita ufficiale sul mercato del set
// ThemeID è l'identificativo del tema a cui appartiene il set
// NumParts è il numero totale di pezzi dichiarati nel set
// PartInventory è la lista contenente i dettagli dei singoli pezzi che compongono il set
// UserReviews è l'elenco delle recensioni lasciate dagli utenti per questo set
// ReviewRating rappresenta la media aggregata delle valutazioni ricevute dagli utenti
type Set struct {
	ID             bson.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	SetNum         string          `bson:"set_num" json:"set_num"`
	Name           string          `bson:"name" json:"name"`
	Year           int             `bson:"year" json:"year"`
	ThemeID        int             `bson:"theme_id" json:"theme_id"`
	NumParts       int             `bson:"num_parts" json:"num_parts"`
	PartsInventory []InventoryPart `bson:"parts_inventory" json:"parts_inventory"`
	UserReviews    []UserReview    `bson:"user_reviews" json:"user_reviews"`
	ReviewRating   float64         `bson:"review_rating" json:"review_rating"`
}
