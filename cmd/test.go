package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mates182/ginshot/formatter"
	"github.com/mates182/ginshot/models" // Update with your actual package path

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [filename]",
	Short: "Reads a JSON file and parses it into the Brewer struct",
	Long: `This command reads a JSON configuration file and parses its content 
into the Brewer struct for validation and testing.`,
	Args: cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		// Read and parse the JSON file
		brewer, err := parseBrewerFile(filename)
		if err != nil {
			fmt.Printf("Error parsing JSON file '%s': %v\n", filename, err)
			return
		}
		fmt.Printf("Successfully parsed JSON file: %s\n", filename)

		if brewer.Root.Gitignore {
			if err := createGitIgnore("./"); err != nil {
				fmt.Println(err)
				return
			}
		}

		/*if brewer.Root.MasterDockerCompose {
			if err := createMasterDockerCompose("./"); err != nil {
				fmt.Println(err)
				return
			}
		}*/

		if brewer.Root.Database.Type != "" {
			addDBConfig(&brewer.Root.Database)
		}

		if brewer.Services != nil {
			for name, service := range brewer.Services {
				config := &models.ProjectConfig{
					ProjectName: name,
					Port:        80,
					Service:     service,
					General:     brewer.General,
					Database:    brewer.Root.Database,
				}

				if err := CreateService(config); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

	},
}

// parseBrewerFile reads and parses a JSON file into the Brewer struct
func parseBrewerFile(filename string) (*models.Brewer, error) {
	// Read the file content
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the JSON content into Brewer struct
	var brewer models.Brewer
	err = json.Unmarshal(data, &brewer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &brewer, nil
}

// CreateServiceDirectory creates the necessary directory structure for a microservice
func CreateService(config *models.ProjectConfig) error {
	fmt.Println(Bold + Cyan + "Generating for: " + config.ProjectName + Reset)
	baseDir := fmt.Sprintf("./%s", config.ProjectName)

	// Define required directories
	directories := []string{
		"cmd",
		"config/cors",
		"internal/controller",
		"internal/service",
		"internal/repository",
		"internal/data/models",
		"internal/data/requests",
		"internal/data/responses",
		"internal/data/messages",
		"internal/event",
		"internal/middleware",
		"internal/secrets",
		"pkg",
		"router",
		"docs",
		"test",
		"deployments",
	}

	// Loop through and create directories
	for _, dir := range directories {
		fullPath := fmt.Sprintf("%s/%s", baseDir, dir)
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	fmt.Printf("Service directory structure created successfully for: %s\n", config.ProjectName)

	// Set up the Go module
	if err := createGoMod(baseDir, config.ProjectName); err != nil {
		fmt.Println(err)
		return err
	}
	// Create Readme File
	if err := createReadmeFile(baseDir, config); err != nil {
		fmt.Println(err)
		return err
	}

	if err := createMainFile(baseDir+"/cmd", config); err != nil {
		fmt.Println(err)
		return err
	}

	if err := createRouterFile(baseDir, config); err != nil {
		fmt.Println(err)
		return err
	}
	// TODO validate
	if err := createCorsFile(baseDir); err != nil {
		fmt.Println(err)
		return err
	}

	if config.Service.Dockerfile || config.General.Dockerfile {
		if err := createDockerfile(baseDir); err != nil {
			fmt.Println(err)
			return err
		}
	}
	if config.Service.DockerCompose || config.General.DockerCompose {
		if err := createDockerCompose(baseDir+"/deployments", config); err != nil {
			fmt.Println(err)
			return err
		}
	}

	generateScaffold(config.General.Models, baseDir+"/internal/data/models/", "model", config, generateModelGo)
	generateScaffold(config.General.Requests, baseDir+"/internal/data/requests/", "request", config, generateRequestGo)
	generateScaffold(config.General.Responses, baseDir+"/internal/data/responses/", "response", config, generateResponseGo)
	generateScaffold(config.General.Messages, baseDir+"/internal/data/messages/", "message", config, generateMessageGo)

	generateScaffold(config.Service.Models, baseDir+"/internal/data/models/", "model", config, generateModelGo)
	generateScaffold(config.Service.Requests, baseDir+"/internal/data/requests/", "request", config, generateRequestGo)
	generateScaffold(config.Service.Responses, baseDir+"/internal/data/responses/", "response", config, generateResponseGo)
	generateScaffold(config.Service.Messages, baseDir+"/internal/data/messages/", "message", config, generateMessageGo)

	generateDBContextFiles(config, config.Database.Type, config.Database.Name, baseDir+"/internal/repository")
	// Pass the name to the function
	generateSecretsFile(config.Database.Type, config.Database.Name, baseDir) // Pass the name to the function

	generateEnvFile(config.Database.Type, config.Database.Name, baseDir)
	addDatabaseToDockerCompose(config.Database.Type, config.Database.Name, baseDir)
	var serviceName, routeName, routeType, requestType, responseType, dbType string
	serviceName = formatter.ToPascalCase(config.ProjectName)

	var crudType int
	switch config.Service.Type {
	case "create":
		crudType = 1
	case "get":
		crudType = 2
	case "update":
		crudType = 3
	case "delete":
		crudType = 4
	case "list":
		crudType = 5
	case "custom":
		crudType = 6
	default:
		fmt.Println(Bold + Yellow + "Invalid selection. Exiting." + Reset)
		return nil
	}

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
		return fmt.Errorf(Bold + Yellow + "Invalid selection. Exiting." + Reset)
	}

	requestType = formatter.ToPascalCase(config.ProjectName) + "Request"
	responseType = formatter.ToPascalCase(config.ProjectName) + "Response"

	dbType = config.Database.Type

	model := config.Database.Model

	id := config.Database.ID

	generateController(config, serviceName, crudType, id, model, baseDir)
	generateService(config.ProjectName, serviceName, requestType, responseType, baseDir, config.Service.Type)
	generateServiceImpl(config, serviceName, crudType, model, id, baseDir, dbType)
	updateRouter(config, routeName, routeType, serviceName, baseDir)

	fmt.Println("Files generated and router updated successfully.")

	// Clean up Go modules
	if err := runGoModTidy(baseDir); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(testCmd)
}
