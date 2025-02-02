package models

type ProjectConfig struct {
	GitIgnore     bool     `json:"GitIgnore" bson:"GitIgnore"`
	Database      Database `json:"Database" bson:"Database"`
	ProjectName   string   `json:"ProjectName" bson:"ProjectName"`
	Port          int      `json:"Port" bson:"Port"`
	Cors          bool     `json:"Cors" bson:"Cors"`
	Dockerfile    bool     `json:"Dockerfile" bson:"Dockerfile"`
	DockerCompose bool     `json:"DockerCompose" bson:"DockerCompose"`
	Service       Service  `json:"Service" bson:"Service"`
	General       Service  `json:"General"`
}
