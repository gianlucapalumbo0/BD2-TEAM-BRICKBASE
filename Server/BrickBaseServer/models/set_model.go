package models

type Set struct {
	SetNum   string `bson:"set_num" json:"set_num" validate:"required"`
	Name     string `bson:"name" json:"name" validate:"required"`
	Year     string `bson:"year" json:"year" validate:"required"`
	ThemeID  string `bson:"theme_id" json:"theme_id" validate:"required"`
	NumParts string `bson:"num_parts" json:"num_parts" validate:"required"`
}
