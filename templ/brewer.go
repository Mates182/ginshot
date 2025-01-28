package templates

import (
	"fmt"

	"github.com/mates182/ginshot/formatter"
	"github.com/mates182/ginshot/models"
)

func GetBrewerTemplate(config *models.ProjectConfig) string {
	projectName := formatter.ToPascalCase(config.ProjectName)
	return fmt.Sprintf(`{
	"project_name": "%s",
	"models": {
		"Data": {
			"Message": "string",
			"Greeting": "Greeting"
		},
		"Greeting": {
			"Hello": "string"
		}
	},
	"requests": {
		"%sRequest": {
			"Data": "models.Data",
			"Message": "string"
		}
	},
	"responses": {
		"%sResponse": {
			"Data": "models.Data",
			"Message": "string"
		}
	}
}`, config.ProjectName, projectName, projectName)
}
