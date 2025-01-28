package templates

import (
	"fmt"

	"github.com/mates182/ginshot/models"
)

func GetReadmeTemplate(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s

	A Gin-based microservice created with ginshot.
	
	## Getting Started
	
	These instructions will help you run the project on your local machine.
	
	### Prerequisites
	
	- Go 1.16 or higher
	
	### Running the service
	
	1. Start the server:
	   `+"```"+`bash
	   go run main.go
	   `+"```"+`
	   `+"```"+`bash
	   curl http://localhost:%d/ping
	   `+"```"+`
	
	## API Endpoints
	
	- GET /ping - Health check endpoint that returns "pong"
	
	## Built With
	
	- [Gin](https://github.com/gin-gonic/gin) - Web framework
	- [ginshot](https://github.com/yourusername/ginshot) - Project scaffolding tool
	
	`, config.ProjectName, config.Port)
}
