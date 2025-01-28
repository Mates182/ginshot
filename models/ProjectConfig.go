package models

type ProjectConfig struct {
	Dockerfile bool `json:"Dockerfile"`
	DockerCompose bool `json:"DockerCompose"`
	GitIgnore bool `json:"GitIgnore"`
	Database Database `json:"Database"`
	ProjectName string `json:"ProjectName"`
	Port int `json:"Port"`
	Cors bool `json:"Cors"`
}
