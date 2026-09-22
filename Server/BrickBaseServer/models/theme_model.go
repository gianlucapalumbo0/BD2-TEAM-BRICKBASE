package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Theme rappresenta un singolo documento all'interno della collezione "themes" in MongoDB
// ID è l'identificativo primario univoco generato da MongoDB
// ThemeID è l'identificativo numerico del tema
// Name è il nome del tema
// ParentID contiene l'ID del tema genitore. Un tema può avere dei sottotemi.
type Theme struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ThemeID  int           `bson:"id" json:"theme_id"`
	Name     string        `bson:"name" json:"name"`
	ParentID *int          `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
}
