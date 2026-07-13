package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Part rappresenta il documento nella collezione "parts"
type Part struct {
	ID      bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PartNum string        `bson:"part_num" json:"part_num"`
	Name    string        `bson:"name" json:"name"`
}
