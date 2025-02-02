package templates

import (
	"fmt"

	"github.com/mates182/ginshot/formatter"
	"github.com/mates182/ginshot/models"
)

func GetControllerTemplate(config *models.ProjectConfig, crudType int, id, model string) string {
	proyectName := formatter.ToPascalCase(config.ProjectName)
	var args string
	if crudType == 2 || crudType == 4 {
		args = GetRequestWithParamsTemplate(id, model, proyectName)
	} else {
		args = GetRequestWithBodyTemplate(proyectName, crudType == 5)
	}

	return fmt.Sprintf(`package controller

import (
	
	"github.com/gin-gonic/gin"
	`+func() string {
		if crudType == 5 {
			return ""
		}
		return `"net/http"
	requests "` + config.ProjectName + `/internal/data/requests"
	responses "` + config.ProjectName + `/internal/data/responses"`
	}()+`
	services "%s/internal/service"
)

type %sController struct {
	%sService services.%sService
}

func New%sController(service services.%sService) *%sController {
	return &%sController{
		%sService: service,
	}
}

func (ctrl *%sController) %s(c *gin.Context) {
	%s
	status, res := ctrl.%sService.%sHandler(`+func() string {
		if crudType == 5 {
			return ""
		}
		return "request"
	}()+`)

	c.IndentedJSON(status, res)
}`, config.ProjectName,
		proyectName, proyectName, proyectName,
		proyectName, proyectName, proyectName, proyectName, proyectName,
		proyectName, proyectName,
		args,
		proyectName, proyectName,
	)
}

func GetRequestWithBodyTemplate(projectName string, isList bool) string {
	if isList {
		return ""
	}
	return fmt.Sprintf(`var request requests.%sRequest
		if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.%sResponse{Message: "Invalid request body"})
		return
		}`, projectName, projectName)

}

func GetRequestWithParamsTemplate(id string, model string, projectName string) string {
	return fmt.Sprintf(`%s := c.Param("id")
	if %s == "" {
		c.IndentedJSON(http.StatusBadRequest, responses.%sResponse{Message: "%s is required"})
		return
	}
	request := requests.%sRequest{%s: %s}
`, id,
		id,
		projectName, id,
		projectName, id, id)
}
