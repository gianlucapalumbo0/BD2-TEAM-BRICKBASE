package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Part rappresenta un singolo documento all'interno della collezione "parts" in MongoDB
// ID è l'identificatore primario univoco generato da MongoDB
// PartNum  è il codice identificativo alfanumerico del pezzo
// Name è il nome ufficiale del pezzo
// PartImgUrl è l'URL dell'immagine del pezzo
type Part struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PartNum    string        `bson:"part_num" json:"part_num"`
	Name       string        `bson:"name" json:"name"`
	PartImgUrl string        `bson:"part_img_url" json:"part_img_url"`
}
