package models

import "go.mongodb.org/mongo-driver/v2/bson"

// InventoryPart rappresenta il singolo pezzo all'interno dell'array del Set
type InventoryPart struct {
	PartNum  string `bson:"part_num" json:"part_num"`
	ColorID  int    `bson:"color_id" json:"color_id"`
	Quantity int    `bson:"quantity" json:"quantity"`
	IsSpare  bool   `bson:"is_spare" json:"is_spare"`
}

type UserReview struct {
	UserID string  `bson:"user_id" json:"user_id"`
	Review string  `bson:"review" json:"review"`
	Rating float64 `bson:"rating" json:"rating"`
}

// Set rappresenta il documento principale nella collezione "sets"
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
