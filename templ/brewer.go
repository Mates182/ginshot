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
		"%s": {
			"ID": "string",
			"Message": "string",
			"Greeting": "Greeting"
		},
		"Greeting": {
			"Hello": "string"
		}
	},
	"requests": {
		"%sRequest": {
			"%s": "models.%s",
			"Message": "string"
		}
	},
	"responses": {
		"%sResponse": {
			"%s": "models.%s",
			"Message": "string"
		}
	}
}`, config.ProjectName,
		config.Database.Model,
		projectName,
		config.Database.Model, config.Database.Model,
		projectName,
		config.Database.Model, config.Database.Model)
}
