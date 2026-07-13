package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Color rappresenta il documento nella collezione "colors"
type Color struct {
	ID      bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ColorID int           `bson:"color_id" json:"color_id"`
	Name    string        `bson:"name" json:"name"`
	RGB     string        `bson:"rgb" json:"rgb"`
	IsTrans bool          `bson:"is_trans" json:"is_trans"`
}
