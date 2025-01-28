package templates

import (
	"fmt"

	"github.com/mates182/ginshot/models"
)

func GetGinshotJSONTemplate(config *models.ProjectConfig) string {
	return fmt.Sprintf(`{
	"ProjectName": "%s",
	"Port": %d,
	"Cors": %t,
	"Dockerfile": %t,
	"DockerCompose": %t,
	"GitIgnore": %t,
	"Services": []
	  }`, config.ProjectName, config.Port, config.Cors, config.Dockerfile, config.DockerCompose, config.GitIgnore)
}
