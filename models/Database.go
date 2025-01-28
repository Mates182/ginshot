package models

type Database struct {
	Table string `json:"Table" bson:"Table"`
	Name string `json:"Name" bson:"Name"`
	Type string `json:"Type" bson:"Type"`
}
