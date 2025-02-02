package models

type Root struct {
	MasterDockerCompose bool     `json:"master_docker_compose" bson:"master_docker_compose"`
	Gitignore           bool     `json:"gitignore" bson:"gitignore"`
	Database            Database `json:"database" bson:"database"`
	Domain              string   `json:"domain" bson:"domain"`
}
