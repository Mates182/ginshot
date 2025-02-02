package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/mates182/ginshot/formatter"
	"github.com/mates182/ginshot/models"
	"github.com/mates182/ginshot/reader"
	templates "github.com/mates182/ginshot/templ"
	"github.com/mates182/ginshot/writer"
	"github.com/spf13/cobra"
)

// Read project name from ginshot.json
func getProjectName() string {
	file, err := os.ReadFile("ginshot.json")
	if err != nil {
		fmt.Println("Error reading ginshot.json:", err)
		os.Exit(1)
	}

	var config models.ProjectConfig
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
		var serviceName, routeName, routeType, requestType, responseType, dbType string
		projectName := getProjectName()
		config, err := reader.LoadConfig()
		if err != nil {
			fmt.Println("Error reading project config:", err)
			return
		}
		serviceName = formatter.ToPascalCase(config.ProjectName)

		var crudType int

		fmt.Print(Bold + Cyan + `Select the CRUD Type:
	` + White + `(1) Create
	` + White + `(2) Read
	` + White + `(3) Update
	` + White + `(4) Delete
	` + White + `(5) List
	` + White + `(6) Custom
` + Grey + `>> ` + Reset)
		fmt.Scanln(&crudType)

		switch crudType {
		case 1:
			routeType = "POST"
			routeName = "/create"
		case 2:
			routeType = "GET"
			routeName = "/get/:id"
		case 3:
			routeType = "PATCH"
			routeName = "/update"
		case 4:
			routeType = "DELETE"
			routeName = "/delete/:id"
		case 5:
			routeType = "GET"
			routeName = "/list"
		case 6:
			fmt.Print("Enter the route name (e.g., /api/ping): ")
			fmt.Scanln(&routeName)

			fmt.Print("Enter the route type (GET, POST, PUT, DELETE): ")
			fmt.Scanln(&routeType)
		default:
			fmt.Println(Bold + Yellow + "Invalid selection. Exiting." + Reset)
			return
		}

		requestType = formatter.ToPascalCase(config.ProjectName) + "Request"
		responseType = formatter.ToPascalCase(config.ProjectName) + "Response"

		dbType = config.Database.Type

		model := config.Database.Model

		id := config.Database.ID

		generateController(config, serviceName, crudType, id, model, "")
		generateService(projectName, serviceName, requestType, responseType, "", "")
		generateServiceImpl(config, serviceName, crudType, model, id, dbType)
		updateRouter(config, routeName, routeType, serviceName, "")

		fmt.Println("Files generated and router updated successfully.")
	},
}

func generateController(config *models.ProjectConfig, serviceName string, crudType int, id, model string, baseDir string) {
	controllerTemplate := templates.GetControllerTemplate(config, crudType, id, model)

	dir := fmt.Sprintf("%s/internal/controller/%s-controller.go", baseDir, formatter.ToLowerCase(serviceName))
	writer.WriteFile(dir, controllerTemplate)
}

func generateService(projectName, serviceName, requestType, responseType string, baseDir string, crudType string) {
	serviceTemplate := `// auto-generated with ginshot
package service

import (
	` + func() string {
		if crudType == "list" {
			return ""
		}
		return `requests "` + projectName + `/internal/data/requests"`
	}() + `
	responses "{{.ProjectName}}/internal/data/responses"
)

type {{.ServiceName}}Service interface {
	{{.ServiceName}}Handler(` + func() string {
		if crudType == "list" {
			return ""
		}
		return `request requests.` + requestType
	}() + `) (int, responses.{{.ResponseType}})
}
`

	fileName := fmt.Sprintf("%s/internal/service/%s-service.go", baseDir, formatter.ToLowerCase(serviceName))
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

func generateServiceImpl(config *models.ProjectConfig, serviceName string, crudType int, model string, id string, baseDir string, db ...string) {
	// Set default value if includeBson is not provided
	dbName := ""
	if len(db) > 0 {
		dbName = db[0]

	}

	serviceImplTemplate := templates.GetServiceImplTemplate(config, dbName, serviceName, model, id, crudType)

	dir := fmt.Sprintf("%s/internal/service/%s-service-impl.go", baseDir, formatter.ToLowerCase(serviceName))
	writer.WriteFile(dir, serviceImplTemplate)
}

func updateRouter(config *models.ProjectConfig, routeName, routeType, serviceName, baseDir string) {
	dbName := config.Database.Name
	routerFileName := baseDir + "/router/router.go"
	file, _ := os.OpenFile(routerFileName, os.O_RDWR, 0644)
	defer file.Close()

	content, _ := os.ReadFile(routerFileName)
	fileContent := string(content)

	insertPoint := "//[ginshot-routes]"
	insertCode := fmt.Sprintf(
		"\n\trouter.%s(\"%s\", controller.New%sController(service.New%sServiceImpl("+func() string {
			if dbName != "" {
				return "dbcontext.GetDBClient()"
			}
			return ""
		}()+")).%s)",
		strings.ToUpper(routeType), routeName, serviceName, serviceName, serviceName,
	)

	if !strings.Contains(fileContent, insertCode) {
		newContent := strings.Replace(fileContent, insertPoint, insertPoint+insertCode, 1)

		// Import statements for both controller and service
		importController := "\n\t\"" + config.ProjectName + "/internal/controller\""
		importService := "\n\t\"" + config.ProjectName + "/internal/service\""
		importDbContext := "\n\t\"" + config.ProjectName + "/internal/repository" + "\""

		importInsertPoint := "import ("
		if !strings.Contains(newContent, importController) {
			newContent = strings.Replace(newContent, importInsertPoint, importInsertPoint+importController, 1)
		}
		if !strings.Contains(newContent, importService) {
			newContent = strings.Replace(newContent, importInsertPoint, importInsertPoint+importService, 1)
		}
		if dbName != "" {
			newContent = strings.Replace(newContent, importInsertPoint, importInsertPoint+importDbContext, 1)
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
