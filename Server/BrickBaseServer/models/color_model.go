package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Color rappresenta un singolo documento all'intero della collezione "colors" in MongoDB.
// ID è l'identificativo primario univoco generato da MongoDB
// ColorID è l'identificativo del colore
// Name è il nome descrittivo del colore
// RGB rappresenta il colore in codice esadecimale
// IsTrans idica se il colore si riferisce ad un materiale trasparente

type Color struct {
	ID      bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ColorID int           `bson:"color_id" json:"color_id"`
	Name    string        `bson:"name" json:"name"`
	RGB     string        `bson:"rgb" json:"rgb"`
	IsTrans bool          `bson:"is_trans" json:"is_trans"`
}
