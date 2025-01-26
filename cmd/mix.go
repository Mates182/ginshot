package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// ProjectConfig represents the structure of ginshot.json
type ProjectConfig struct {
	ProjectName string `json:"project_name"`
}

// Read project name from ginshot.json
func getProjectName() string {
	file, err := os.ReadFile("ginshot.json")
	if err != nil {
		fmt.Println("Error reading ginshot.json:", err)
		os.Exit(1)
	}

	var config ProjectConfig
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Println("Error parsing ginshot.json:", err)
		os.Exit(1)
	}
	return config.ProjectName
}

// mixCmd represents the mix command
var mixCmd = &cobra.Command{
	Use:   "mix",
	Short: "Generate controller, service, and route files",
	Long:  "This command generates files and updates the router based on user inputs and project configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		var serviceName, routeName, routeType, requestType, responseType string
		projectName := getProjectName()

		fmt.Print("Enter the service name: ")
		fmt.Scanln(&serviceName)

		fmt.Print("Enter the route name (e.g., /api/ping): ")
		fmt.Scanln(&routeName)

		fmt.Print("Enter the route type (GET, POST, PUT, DELETE): ")
		fmt.Scanln(&routeType)

		fmt.Print("Enter the request type (e.g., PingRequest): ")
		fmt.Scanln(&requestType)

		fmt.Print("Enter the response type (e.g., PingResponse): ")
		fmt.Scanln(&responseType)

		generateController(projectName, serviceName, requestType, responseType)
		generateService(projectName, serviceName, requestType, responseType)
		generateServiceImpl(projectName, serviceName, requestType, responseType)
		updateRouter(projectName, routeName, routeType, serviceName)

		fmt.Println("Files generated and router updated successfully.")
	},
}

func generateController(projectName, serviceName, requestType, responseType string) {
	controllerTemplate := `// auto-generated with ginshot
package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	requests "{{.ProjectName}}/data/requests"
	responses "{{.ProjectName}}/data/responses"
	services "{{.ProjectName}}/service"
)

type {{.ServiceName}}Controller struct {
	{{.ServiceName}}Service services.{{.ServiceName}}Service
}

func New{{.ServiceName}}Controller(service services.{{.ServiceName}}Service) *{{.ServiceName}}Controller {
	return &{{.ServiceName}}Controller{
		{{.ServiceName}}Service: service,
	}
}

func (ctrl *{{.ServiceName}}Controller) {{.ServiceName}}(c *gin.Context) {
	var request requests.{{.RequestType}}
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.{{.ResponseType}}{Message: "Invalid request body"})
		return
	}
	status, res := ctrl.{{.ServiceName}}Service.{{.ServiceName}}Handler(request)

	c.IndentedJSON(status, res)
}
`

	fileName := fmt.Sprintf("./controller/%s-controller.go", serviceName)
	file, _ := os.Create(fileName)
	defer file.Close()

	tmpl, _ := template.New("controller").Parse(controllerTemplate)
	tmpl.Execute(file, map[string]string{
		"ServiceName":  serviceName,
		"RequestType":  requestType,
		"ResponseType": responseType,
		"ProjectName":  projectName,
	})
}

func generateService(projectName, serviceName, requestType, responseType string) {
	serviceTemplate := `// auto-generated with ginshot
package service

import (
	requests "{{.ProjectName}}/data/requests"
	responses "{{.ProjectName}}/data/responses"
)

type {{.ServiceName}}Service interface {
	{{.ServiceName}}Handler(request requests.{{.RequestType}}) (int, responses.{{.ResponseType}})
}
`

	fileName := fmt.Sprintf("./service/%s-service.go", serviceName)
	file, _ := os.Create(fileName)
	defer file.Close()

	tmpl, _ := template.New("service").Parse(serviceTemplate)
	tmpl.Execute(file, map[string]string{
		"ProjectName":  projectName,
		"ServiceName":  serviceName,
		"RequestType":  requestType,
		"ResponseType": responseType,
	})
}

func generateServiceImpl(projectName, serviceName, requestType, responseType string) {
	serviceImplTemplate := `// auto-generated with ginshot
package service

import (
	requests "{{.ProjectName}}/data/requests"
	responses "{{.ProjectName}}/data/responses"
	"net/http"
)
type {{.ServiceName}}ServiceImpl struct {
	// Add Components
}

func New{{.ServiceName}}ServiceImpl() {{.ServiceName}}Service {
	return &{{.ServiceName}}ServiceImpl{
		// Add Components
	}
}

func (service *{{.ServiceName}}ServiceImpl) {{.ServiceName}}Handler(request requests.{{.RequestType}}) (int, responses.{{.ResponseType}}) {
	response := responses.{{.ResponseType}}{}
	return http.StatusOK, response
}
`

	fileName := fmt.Sprintf("./service/%s-service-impl.go", serviceName)
	file, _ := os.Create(fileName)
	defer file.Close()

	tmpl, _ := template.New("serviceImpl").Parse(serviceImplTemplate)
	tmpl.Execute(file, map[string]string{
		"ProjectName":  projectName,
		"ServiceName":  serviceName,
		"RequestType":  requestType,
		"ResponseType": responseType,
	})
}

func updateRouter(projectName, routeName, routeType, serviceName string) {
	routerFileName := "./router/router.go"
	file, _ := os.OpenFile(routerFileName, os.O_RDWR, 0644)
	defer file.Close()

	content, _ := os.ReadFile(routerFileName)
	fileContent := string(content)

	insertPoint := "//[ginshot-routes]"
	insertCode := fmt.Sprintf(
		"\n\trouter.%s(\"%s\", controller.New%sController(service.New%sServiceImpl()).%s)",
		strings.ToUpper(routeType), routeName, serviceName, serviceName, serviceName,
	)

	if !strings.Contains(fileContent, insertCode) {
		newContent := strings.Replace(fileContent, insertPoint, insertPoint+insertCode, 1)

		// Import statements for both controller and service
		importController := "\n\t\"" + projectName + "/controller\""
		importService := "\n\t\"" + projectName + "/service\""

		importInsertPoint := "import ("
		if !strings.Contains(newContent, importController) {
			newContent = strings.Replace(newContent, importInsertPoint, importInsertPoint+importController, 1)
		}
		if !strings.Contains(newContent, importService) {
			newContent = strings.Replace(newContent, importInsertPoint, importInsertPoint+importService, 1)
		}

		os.WriteFile(routerFileName, []byte(newContent), 0644)
		fmt.Println("Route and imports added successfully.")
	} else {
		fmt.Println("Route already exists, skipping addition.")
	}
}

func init() {
	rootCmd.AddCommand(mixCmd)
}
