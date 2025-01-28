package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/mates182/ginshot/formatter"
	"github.com/mates182/ginshot/models"
	templates "github.com/mates182/ginshot/templ"
	"github.com/mates182/ginshot/writer"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Gin project",
	Long: `Initialize a new Gin-based microservice with a structured project layout.
For example:

ginshot init my-service`,
	Args: cobra.MaximumNArgs(1), // Allows zero or one argument
	Run:  generateProject,
}

var port string

func init() {
	// Add the port flag to the init command
	initCmd.Flags().StringVarP(&port, "port", "p", "", "Port for the microservice (default is 8080)")

	rootCmd.AddCommand(initCmd)
}

// generateProject creates the project structure based on the provided name or asks for it
func generateProject(cmd *cobra.Command, args []string) {
	portString, err := strconv.Atoi(port)
	if err != nil {
		portString = 8080
	}
	config := &models.ProjectConfig{
		ProjectName:   "",
		Port:          portString,
		Cors:          true,
		Dockerfile:    true,
		DockerCompose: true,
		GitIgnore:     true,
		Database: models.Database{
			Model: "Data",
			ID:    "ID",
		},
	}

	// Prompt for project name if it's not provided in the arguments
	if len(args) == 0 {
		fmt.Print(Bold + Cyan + "Project name " + Magenta + "(format: lower-case-project-name) " + Grey + ">> " + Reset)
		fmt.Scanln(&config.ProjectName)
	} else {
		config.ProjectName = args[0]
	}

	// Validate the project name
	if regexp.MustCompile(`^[a-zA-Z][\da-zA-Z]*([-:\/\\_*+.\[\]{}()|"',;<>?]+[\da-zA-Z]+)*$`).FindString(config.ProjectName) == "" {
		fmt.Println(Red + "Error: Project name cannot be empty, have blank spaces nor start with numbers or special characters" + Reset)
		return
	}

	// Ask for the port if not provided by the flag
	if port == "" {
		fmt.Print(Bold + Cyan + "Port " + Magenta + "(default is 8080) " + Grey + ">> " + Reset)
		fmt.Scanln(&port)
	}

	// Default to port 8080 if the user doesn't provide a port
	if regexp.MustCompile(`^[\d]{1,5}$`).FindString(fmt.Sprintf("%d", config.Port)) == "" {
		fmt.Println(Yellow + "Using default port: 8080" + Reset)
	} else {

	}

	config.Port, err = strconv.Atoi(port)
	if err != nil {
		fmt.Println(Yellow + "Using default port: 8080" + Reset)
		config.Port = 8080
	}

	fmt.Print(Bold + Cyan + "Model Name:" + Magenta + "(recomended format: PascalCase) " + Grey + ">> " + Reset)
	fmt.Scanln(&config.Database.Model)

	fmt.Print(Bold + Cyan + "Model ID Name:" + Magenta + "(recomended format: PascalCase) " + Grey + ">> " + Reset)
	fmt.Scanln(&config.Database.ID)

	if err := generateProjectFiles(config); err != nil {
		fmt.Println(err)
		return
	}

	// Successfully created the microservice
	fmt.Println(Bold + "Project '" + Cyan + config.ProjectName + White + "' created " + Green + "successfully" + White + " in '" + Bold + Cyan + "./" + config.ProjectName + White + "' on port " + Bold + Cyan + strconv.Itoa(config.Port) + Reset)
	// printProjectInstructions prints instructions for running the project

	fmt.Println(Bold + "\nTo run your project:" + Reset)
	fmt.Printf(Grey+"  cd %s\n", formatter.ToPascalCase(config.ProjectName))
	fmt.Println("  go run main.go" + Reset)
	fmt.Println(Bold + "\nTest the API:" + Reset)
	fmt.Printf(Grey+"  curl http://localhost:%d/ping\n"+Reset, config.Port)
}

func generateProjectFiles(config *models.ProjectConfig) error {
	// Create the project structure
	projectName := formatter.ToPascalCase(config.ProjectName)
	dir := fmt.Sprintf("./%s", projectName)
	if err := createProjectDirectory(dir); err != nil {
		fmt.Println(err)
		return err
	}

	// Create the necessary files
	if err := createMainFile(dir, config); err != nil {
		fmt.Println(err)
		return err
	}

	if err := createRouterFile(dir, config); err != nil {
		fmt.Println(err)
		return err
	}
	if err := createCorsFile(dir); err != nil {
		fmt.Println(err)
		return err
	}

	if err := createDockerfile(dir); err != nil {
		fmt.Println(err)
		return err
	}

	if err := createDockerCompose(dir, config); err != nil {
		fmt.Println(err)
		return err
	}

	if err := createGitIgnore(dir); err != nil {
		fmt.Println(err)
		return err
	}
	if err := createGinshotJSON(dir, config); err != nil {
		fmt.Println(err)
		return err
	}

	// Set up the Go module
	if err := createGoMod(dir, config.ProjectName); err != nil {
		fmt.Println(err)
		return err
	}

	// Clean up Go modules
	if err := runGoModTidy(dir); err != nil {
		fmt.Println(err)
		return err
	}
	// Create Readme File
	if err := createReadmeFile(dir, config); err != nil {
		fmt.Println(err)
		return err
	}

	// Create Brewer Template
	if err := createBrewerTemplateJSON(dir, config); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// createProjectDirectory creates the project directory and necessary subdirectories
func createProjectDirectory(dir string) error {
	// Create the main project directory
	if err := writer.CreateDir(dir); err != nil {
		return err
	}

	// Create the 'router' directory
	if err := writer.CreateDir(fmt.Sprintf("%s/router", dir)); err != nil {
		return err
	}

	// Create the 'config' directory
	if err := writer.CreateDir(fmt.Sprintf("%s/config", dir)); err != nil {
		return err
	}

	// Create the 'data' subdirectories
	if err := writer.CreateDir(fmt.Sprintf("%s/data/requests", dir)); err != nil {
		return err
	}
	if err := writer.CreateDir(fmt.Sprintf("%s/data/responses", dir)); err != nil {
		return err
	}

	// Create the 'models' directory
	if err := writer.CreateDir(fmt.Sprintf("%s/models", dir)); err != nil {
		return err
	}

	// Create the 'brewer' directory
	if err := writer.CreateDir(fmt.Sprintf("%s/brewer", dir)); err != nil {
		return err
	}

	// Create the 'controller' directory
	if err := writer.CreateDir(fmt.Sprintf("%s/controller", dir)); err != nil {
		return err
	}

	// Create the 'service' directory
	if err := writer.CreateDir(fmt.Sprintf("%s/service", dir)); err != nil {
		return err
	}

	return nil
}

// createMainFile creates the main.go file
func createMainFile(dir string, config *models.ProjectConfig) error {
	mainFile := templates.GetMainTemplate(config)
	return writer.WriteFile(dir+"/main.go", mainFile)
}

// createRouterFile creates the router.go file
func createRouterFile(dir string, config *models.ProjectConfig) error {
	routerFile := templates.GetRouterTemplate(config)
	return writer.WriteFile(dir+"/router/router.go", routerFile)
}

// createGoMod initializes the Go module
func createGoMod(dir, name string) error {
	cmdGoMod := exec.Command("go", "mod", "init", name)
	cmdGoMod.Dir = dir
	if err := cmdGoMod.Run(); err != nil {
		return fmt.Errorf("error creating go.mod: %v", err)
	}
	return nil
}

// runGoModTidy runs `go mod tidy` to clean up the dependencies
func runGoModTidy(dir string) error {
	cmdGoModTidy := exec.Command("go", "mod", "tidy")
	cmdGoModTidy.Dir = dir
	if err := cmdGoModTidy.Run(); err != nil {
		return fmt.Errorf("error running go mod tidy: %v", err)
	}
	return nil
}

// createReadmeFile creates a README.md file with project documentation
func createReadmeFile(dir string, config *models.ProjectConfig) error {
	readmeContent := templates.GetReadmeTemplate(config)
	return writer.WriteFile(dir+"/README.md", readmeContent)
}

// createCorsFile creates the cors.go file with CORS configuration
func createCorsFile(dir string) error {
	corsFile := templates.GetCORSTemplate()
	if err := os.MkdirAll(fmt.Sprintf("%s/config/cors", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating config/cors directory: %v", err)
	}
	return writer.WriteFile(dir+"/config/cors/cors.go", corsFile)
}

// createDockerfile creates a Dockerfile for the project
func createDockerfile(dir string) error {
	dockerfileContent := templates.GetDockerfileTemplate()
	if err := writer.WriteFile(dir+"/Dockerfile", dockerfileContent); err != nil {
		return err
	}

	dockerignoreContent := templates.GetDockerIgnoreTemplate()
	return writer.WriteFile(dir+"/.dockerignore", dockerignoreContent)
}

// createDockerCompose creates a docker-compose.yml file for the project
func createDockerCompose(dir string, config *models.ProjectConfig) error {
	dockerComposeContent := templates.GetDockerComposeTemplate(config)
	return writer.WriteFile(dir+"/docker-compose.yml", dockerComposeContent)
}

// createGitIgnore creates a .gitignore file for the project
func createGitIgnore(dir string) error {
	gitignoreContent := templates.GetGitIgnoreTemplate()
	return writer.WriteFile(dir+"/.gitignore", gitignoreContent)
}

// createGinshotJSON creates the ginshot.json file for the project
func createGinshotJSON(dir string, config *models.ProjectConfig) error {
	ginshotJSONContent := templates.GetGinshotJSONTemplate(config)
	return writer.WriteFile(dir+"/ginshot.json", ginshotJSONContent)
}

// createBrewerTemplateJSON creates the brewer/template.json file
func createBrewerTemplateJSON(dir string, config *models.ProjectConfig) error {
	dataJSONContent := templates.GetBrewerTemplate(config)
	return writer.WriteFile(dir+"/brewer/template.json", dataJSONContent)
}
