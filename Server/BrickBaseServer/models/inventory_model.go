package models

type InventoryPart struct {
	PartNum  string `bson:"part_num" json:"part_num"`
	ColorID  int    `bson:"color_id" json:"color_id"`
	Quantity int    `bson:"quantity" json:"quantity"`
	IsSpare  bool   `bson:"is_spare" json:"is_spare"`
}

type InventorySet struct {
	SetNum   string `bson:"set_num" json:"set_num"`
	Quantity int    `bson:"quantity" json:"quantity"`
}

type Inventory struct {
	ID      int             `bson:"_id" json:"id"`
	Version int             `bson:"version" json:"version"`
	SetNum  string          `bson:"set_num" json:"set_num"`
	Parts   []InventoryPart `bson:"parts,omitempty" json:"parts,omitempty"`
	Sets    []InventorySet  `bson:"sets,omitempty" json:"sets,omitempty"`
}
