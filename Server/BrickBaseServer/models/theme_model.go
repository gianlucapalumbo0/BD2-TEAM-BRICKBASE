package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Theme rappresenta il documento nella collezione "themes"
type Theme struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ThemeID  int           `bson:"theme_id" json:"theme_id"`
	Name     string        `bson:"name" json:"name"`
	ParentID *int          `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
}
