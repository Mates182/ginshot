package models

type Database struct {
	Model      string `json:"model" bson:"Model"`
	ID         string `json:"id" bson:"ID"`
	Name       string `json:"name" bson:"Name"`
	Type       string `json:"type" bson:"Type"`
	Collection string `json:"collection" bson:"Table"`
}
