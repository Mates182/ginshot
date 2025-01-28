package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/mates182/ginshot/models"
	templates "github.com/mates182/ginshot/templ"
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

type InitConfig struct {
	Name          string
	Port          string
	Cors          bool
	Dockerfile    bool
	DockerCompose bool
	GitIgnore     bool
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
		Dockerfile:    false,
		DockerCompose: false,
		GitIgnore:     true,
	}

	// Prompt for project name if it's not provided in the arguments
	if len(args) == 0 {
		fmt.Print(Bold + Cyan + "Project name " + Magenta + "(recomended format: lower-case-project-name) " + Grey + ">> " + Reset)
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
		config.Port, _ = strconv.Atoi(port)

	}

	if err := generateProjectFiles(config); err != nil {
		fmt.Println(err)
		return
	}

	// Successfully created the microservice
	fmt.Println(Bold + "Project '" + Cyan + config.ProjectName + White + "' created " + Green + "successfully" + White + " in '" + Bold + Cyan + "./" + config.ProjectName + White + "' on port " + Bold + Cyan + strconv.Itoa(config.Port) + Reset)
	// printProjectInstructions prints instructions for running the project

	fmt.Println(Bold + "\nTo run your project:" + Reset)
	fmt.Printf(Grey+"  cd %s\n", config.ProjectName)
	fmt.Println("  go run main.go" + Reset)
	fmt.Println(Bold + "\nTest the API:" + Reset)
	fmt.Printf(Grey+"  curl http://localhost:%d/ping\n"+Reset, config.Port)
}

func generateProjectFiles(config *models.ProjectConfig) error {
	// Create the project structure
	dir := fmt.Sprintf("./%s", config.ProjectName)
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

// createProjectDirectory creates the project directory
func createProjectDirectory(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating project directory: %v", err)
	}
	// Create the 'router' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/router", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating router directory: %v", err)
	}
	// Create the 'config' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/config", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}
	// Create the 'data' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/data/requests", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/data/responses", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/models", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}
	// Create the 'data' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/brewer", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/controller", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/service", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}
	return nil

}

// createMainFile creates the main.go file
func createMainFile(dir string, config *models.ProjectConfig) error {
	mainFile := templates.GetMainTemplate(config)

	if err := os.WriteFile(dir+"/main.go", []byte(mainFile), 0644); err != nil {
		return fmt.Errorf("error creating main.go: %v", err)
	}
	return nil
}

// createRouterFile creates the router.go file
func createRouterFile(dir string, config *models.ProjectConfig) error {
	routerFile := templates.GetRouterTemplate(config)

	if err := os.WriteFile(dir+"/router/router.go", []byte(routerFile), 0644); err != nil {
		return fmt.Errorf("error creating router.go: %v", err)
	}
	return nil
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

	if err := os.WriteFile(dir+"/README.md", []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("error creating README.md: %v", err)
	}
	return nil
}

// createCorsFile creates the cors.go file with CORS configuration
func createCorsFile(dir string) error {
	corsFile := templates.GetCORSTemplate()
	if err := os.MkdirAll(fmt.Sprintf("%s/config/cors", dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating config/cors directory: %v", err)
	}

	if err := os.WriteFile(dir+"/config/cors/cors.go", []byte(corsFile), 0644); err != nil {
		return fmt.Errorf("error creating cors.go: %v", err)
	}
	return nil
}

// createDockerfile creates a Dockerfile for the project
func createDockerfile(dir string) error {
	dockerfileContent := templates.GetDockerfileTemplate()

	if err := os.WriteFile(dir+"/Dockerfile", []byte(dockerfileContent), 0644); err != nil {
		return fmt.Errorf("error creating Dockerfile: %v", err)
	}

	dockerignoreContent := templates.GetDockerIgnoreTemplate()

	if err := os.WriteFile(dir+"/.dockerignore", []byte(dockerignoreContent), 0644); err != nil {
		return fmt.Errorf("error creating .dockerignore: %v", err)
	}

	return nil
}

// createDockerCompose creates a docker-compose.yml file for the project
func createDockerCompose(dir string, config *models.ProjectConfig) error {
	dockerComposeContent := templates.GetDockerComposeTemplate(config)

	if err := os.WriteFile(dir+"/docker-compose.yml", []byte(dockerComposeContent), 0644); err != nil {
		return fmt.Errorf("error creating docker-compose.yml: %v", err)
	}

	return nil
}

func createGitIgnore(dir string) error {
	gitignoreContent := templates.GetGitIgnoreTemplate()

	if err := os.WriteFile(dir+"/.gitignore", []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("error creating .gitignore: %v", err)
	}

	return nil
}

func createGinshotJSON(dir string, config *models.ProjectConfig) error {
	ginshotJSONContent := templates.GetGinshotJSONTemplate(config)

	if err := os.WriteFile(dir+"/ginshot.json", []byte(ginshotJSONContent), 0644); err != nil {
		return fmt.Errorf("error creating ginshot.json: %v", err)
	}

	return nil
}

func createBrewerTemplateJSON(dir string, config *models.ProjectConfig) error {
	dataJSONContent := templates.GetBrewerTemplate(config)

	if err := os.WriteFile(dir+"/brewer/template.json", []byte(dataJSONContent), 0644); err != nil {
		return fmt.Errorf("error creating template.json: %v", err)
	}

	return nil
}
