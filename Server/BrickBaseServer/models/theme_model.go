package models

type Theme struct {
	ID       int    `bson:"_id" json:"id"`
	Name     string `bson:"name" json:"name"`
	ParentID *int   `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
}
