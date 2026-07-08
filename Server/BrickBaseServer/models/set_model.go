package models

type Set struct {
	SetNum   string `bson:"set_num" json:"set_num"`
	Name     string `bson:"name" json:"name"`
	Year     string `bson:"year" json:"year"`
	ThemeID  string `bson:"theme_id" json:"theme_id"`
	NumParts string `bson:"num_parts" json:"num_parts"`
}
