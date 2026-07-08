package models

type Color struct {
	ID      int    `bson:"_id" json:"id"`
	Name    string `bson:"name" json:"name"`
	RGB     string `bson:"rgb" json:"rgb"`
	IsTrans bool   `bson:"is_trans" json:"is_trans"`
}
