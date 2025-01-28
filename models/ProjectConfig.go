package models

type ProjectConfig struct {
	Cors bool `json:"Cors" bson:"Cors"`
	Dockerfile bool `json:"Dockerfile" bson:"Dockerfile"`
	DockerCompose bool `json:"DockerCompose" bson:"DockerCompose"`
	GitIgnore bool `json:"GitIgnore" bson:"GitIgnore"`
	Services []Service `json:"Services" bson:"Services"`
	ProjectName string `json:"ProjectName" bson:"ProjectName"`
	Port int `json:"Port" bson:"Port"`
}
