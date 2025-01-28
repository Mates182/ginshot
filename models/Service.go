package models

type Service struct {
	Name string `json:"Name" bson:"Name"`
	Request string `json:"Request" bson:"Request"`
	Response string `json:"Response" bson:"Response"`
}
