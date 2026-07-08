package models

type PartCategory struct {
	ID   int    `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Part struct {
	PartNum   string `bson:"_id" json:"part_num"`
	Name      string `bson:"name" json:"name"`
	PartCatID int    `bson:"part_cat_id" json:"part_cat_id"`
}
