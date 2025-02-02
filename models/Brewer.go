package models

type Brewer struct {
	Root     Root                              `json:"root"`
	General  Service                           `json:"general"`
	Models   map[string]map[string]interface{} `json:"models"`
	Services map[string]Service                `json:"services"`
}
type Service struct {
	Type          string                            `json:"type"`
	Port          int                               `json:"port"`
	Dockerfile    bool                              `json:"dockerfile"`
	DockerCompose bool                              `json:"docker_compose"`
	Models        map[string]map[string]interface{} `json:"models"`
	Requests      map[string]map[string]interface{} `json:"requests"`
	Responses     map[string]map[string]interface{} `json:"responses"`
	Messages      map[string]map[string]interface{} `json:"messages"`
	Cors          Cors                              `json:"cors"`
}

type Cors struct {
	Methods []string `json:"methods"`
}
