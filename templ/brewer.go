package templates

import (
	"fmt"

	"github.com/mates182/ginshot/models"
)

func GetBrewerTemplate(config *models.ProjectConfig) string {
	return fmt.Sprintf(`{
	"project_name": "%s",
	"models": {
		"Ping": {
			"Message": "string",
			"Greeting": "Greeting"
		},
		"Greeting": {
			"Hello": "string"
		}
	},
	"requests": {
		"PingRequest": {
			"Data": "models.Ping",
			"Message": "string"
		}
	},
	"responses": {
		"PingResponse": {
			"Data": "models.Ping",
			"Message": "string"
		}
	}
}`, config.ProjectName)
}
