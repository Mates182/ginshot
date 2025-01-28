package templates

import (
	"fmt"

	"github.com/mates182/ginshot/models"
)

func GetRouterTemplate(config *models.ProjectConfig) string {
	return fmt.Sprintf(`package router

import (
	"github.com/gin-gonic/gin"
	"%s/config/cors"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.GetCORSConfig())
	//[ginshot-routes]
	//[HttpGET] Ping to %s API
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
	`, config.ProjectName, config.ProjectName)
}
