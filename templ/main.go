package templates

import (
	"fmt"

	"github.com/mates182/ginshot/models"
)

func GetMainTemplate(config *models.ProjectConfig) string {
	return fmt.Sprintf(`package main

import (
	"%s/router"
	"fmt"
)

func main() {
	fmt.Println("%s API started!")
	router := router.SetupRouter()
	router.Run("0.0.0.0:%d")
}
		`, config.ProjectName, config.ProjectName, config.Port)
}
