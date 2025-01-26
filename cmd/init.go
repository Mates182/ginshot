package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

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
	config := &InitConfig{
		Name:          "",
		Port:          port,
		Cors:          true,
		Dockerfile:    false,
		DockerCompose: false,
		GitIgnore:     true,
	}

	// Prompt for project name if it's not provided in the arguments
	if len(args) == 0 {
		fmt.Print(Bold + Cyan + "Project name " + Magenta + "(recomended format: lower-case-project-name) " + Grey + ">> " + Reset)
		fmt.Scanln(&config.Name)
	} else {
		config.Name = args[0]
	}

	// Validate the project name
	if regexp.MustCompile(`^[a-zA-Z][\da-zA-Z]*([-:\/\\_*+.\[\]{}()|"',;<>?]+[\da-zA-Z]+)*$`).FindString(config.Name) == "" {
		fmt.Println(Red + "Error: Project name cannot be empty, have blank spaces nor start with numbers or special characters" + Reset)
		return
	}

	// Ask for the port if not provided by the flag
	if config.Port == "" {
		fmt.Print(Bold + Cyan + "Port " + Magenta + "(default is 8080) " + Grey + ">> " + Reset)
		fmt.Scanln(&config.Port)
	}

	// Default to port 8080 if the user doesn't provide a port
	if regexp.MustCompile(`^[\d]{1,5}$`).FindString(config.Port) == "" {
		fmt.Println(Yellow + "Using default port: 8080" + Reset)
		config.Port = "8080"
	}

	if err := generateProjectFiles(config); err != nil {
		fmt.Println(err)
		return
	}

	// Successfully created the microservice
	fmt.Println(Bold + "Project '" + Cyan + config.Name + White + "' created " + Green + "successfully" + White + " in '" + Bold + Cyan + "./" + config.Name + White + "' on port " + Bold + Cyan + config.Port + Reset)
	// printProjectInstructions prints instructions for running the project

	fmt.Println(Bold + "\nTo run your project:" + Reset)
	fmt.Printf(Grey+"  cd %s\n", config.Name)
	fmt.Println("  go run main.go" + Reset)
	fmt.Println(Bold + "\nTest the API:" + Reset)
	fmt.Printf(Grey+"  curl http://localhost:%s/ping\n"+Reset, config.Port)
}

func generateProjectFiles(config *InitConfig) error {
	// Create the project structure
	dir := fmt.Sprintf("./%s", config.Name)
	if err := createProjectDirectory(dir); err != nil {
		fmt.Println(err)
		return err
	}

	// Create the necessary files
	if err := createMainFile(dir, config.Name, config.Port); err != nil {
		fmt.Println(err)
		return err
	}

	if err := createRouterFile(dir, config.Name); err != nil {
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
	if err := createGoMod(dir, config.Name); err != nil {
		fmt.Println(err)
		return err
	}

	// Clean up Go modules
	if err := runGoModTidy(dir); err != nil {
		fmt.Println(err)
		return err
	}
	// Create Readme File
	if err := createReadmeFile(dir, config.Name, config.Port); err != nil {
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
		return fmt.Errorf("Error creating project directory: %v", err)
	}
	// Create the 'router' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/router", dir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating router directory: %v", err)
	}
	// Create the 'config' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/config", dir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating config directory: %v", err)
	}
	// Create the 'data' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/data/requests", dir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating data directory: %v", err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/data/responses", dir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating data directory: %v", err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/models", dir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating data directory: %v", err)
	}
	// Create the 'data' directory
	if err := os.MkdirAll(fmt.Sprintf("%s/brewer", dir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating data directory: %v", err)
	}
	return nil

}

// createMainFile creates the main.go file
func createMainFile(dir, name, port string) error {
	mainFile := fmt.Sprintf(`// auto-generated by ginshot
package main

import (
	"%s/router"
	"fmt"
)

func main() {
	fmt.Println("%s API started!")
	router := router.SetupRouter()
	router.Run("0.0.0.0:%s")
}
	`, name, name, port)

	if err := os.WriteFile(dir+"/main.go", []byte(mainFile), 0644); err != nil {
		return fmt.Errorf("Error creating main.go: %v", err)
	}
	return nil
}

// createRouterFile creates the router.go file
func createRouterFile(dir, name string) error {
	routerFile := fmt.Sprintf(`//auto-generated by ginshot
package router

import (
	"github.com/gin-gonic/gin"
	"`+name+`/config/cors"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.GetCORSConfig())
	// [HttpGET] Ping to %s API
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}
	`, name)

	if err := os.WriteFile(dir+"/router/router.go", []byte(routerFile), 0644); err != nil {
		return fmt.Errorf("Error creating router.go: %v", err)
	}
	return nil
}

// createGoMod initializes the Go module
func createGoMod(dir, name string) error {
	cmdGoMod := exec.Command("go", "mod", "init", name)
	cmdGoMod.Dir = dir
	if err := cmdGoMod.Run(); err != nil {
		return fmt.Errorf("Error creating go.mod: %v", err)
	}
	return nil
}

// runGoModTidy runs `go mod tidy` to clean up the dependencies
func runGoModTidy(dir string) error {
	cmdGoModTidy := exec.Command("go", "mod", "tidy")
	cmdGoModTidy.Dir = dir
	if err := cmdGoModTidy.Run(); err != nil {
		return fmt.Errorf("Error running go mod tidy: %v", err)
	}
	return nil
}

// createReadmeFile creates a README.md file with project documentation
func createReadmeFile(dir, name, port string) error {
	readmeContent := fmt.Sprintf(`# %s

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
   curl http://localhost:%s/ping
   `+"```"+`

## API Endpoints

- GET /ping - Health check endpoint that returns "pong"

## Built With

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [ginshot](https://github.com/yourusername/ginshot) - Project scaffolding tool

`, name, port)

	if err := os.WriteFile(dir+"/README.md", []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("Error creating README.md: %v", err)
	}
	return nil
}

// createCorsFile creates the cors.go file with CORS configuration
func createCorsFile(dir string) error {
	corsFile := `// auto-generated by ginshot
package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetCORSConfig() gin.HandlerFunc {
	corsConfig := cors.New(cors.Config{
		// Set to true to allow all origins (remove if you want to allow specific origins only)
		AllowAllOrigins: true, 

		// Uncomment and modify the line below to allow specific origins instead of all
		// AllowOrigins: []string{"http://localhost:80", "https://example.com"}, 

		// Define allowed HTTP methods (adjust according to your API needs)
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},

		// Specify the allowed headers (remove or add headers as required by your application)
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},

		// Set to true to allow credentials such as cookies or authorization headers
		AllowCredentials: true,
	})

	return corsConfig
}

`
	if err := os.MkdirAll(fmt.Sprintf("%s/config/cors", dir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating config/cors directory: %v", err)
	}

	if err := os.WriteFile(dir+"/config/cors/cors.go", []byte(corsFile), 0644); err != nil {
		return fmt.Errorf("Error creating cors.go: %v", err)
	}
	return nil
}

// createDockerfile creates a Dockerfile for the project
func createDockerfile(dir string) error {
	dockerfileContent := `# auto-generated by ginshot
# change the version if needed
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]
`

	if err := os.WriteFile(dir+"/Dockerfile", []byte(dockerfileContent), 0644); err != nil {
		return fmt.Errorf("Error creating Dockerfile: %v", err)
	}

	dockerignoreContent := `# auto-generated by ginshot
# Ginshot Files (optional)
ginshot.json
brewer


# Binaries and build artifacts
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# IDE files
.idea/
.vscode/
*.swp
*.swo

# Dependencies
vendor/

# Git
.git
.gitignore

# Logs
*.log

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Secrets
.env
.env.local
.env.*.local
`

	if err := os.WriteFile(dir+"/.dockerignore", []byte(dockerignoreContent), 0644); err != nil {
		return fmt.Errorf("Error creating .dockerignore: %v", err)
	}

	return nil
}

// createDockerCompose creates a docker-compose.yml file for the project
func createDockerCompose(dir string, config *InitConfig) error {
	dockerComposeContent := `# auto-generated by ginshot
# Uncomment if needed
# version: '3.8'

services:
  ` + config.Name + `:
    build: .
    ports:
      - "` + config.Port + `:` + config.Port + `"
    container_name: ` + config.Name + `
    environment:
      - GIN_MODE=release
    volumes:
      - .:/app
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge`

	if err := os.WriteFile(dir+"/docker-compose.yml", []byte(dockerComposeContent), 0644); err != nil {
		return fmt.Errorf("Error creating docker-compose.yml: %v", err)
	}

	return nil
}

func createGitIgnore(dir string) error {
	gitignoreContent := `# auto-generated by ginshot
# Ginshot Files (optional)
# ginshot.json
# brewer

# Binaries and build artifacts
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# IDE files
.idea/
.vscode/
*.swp
*.swo

# Dependencies
vendor/

# Logs
*.log

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Secrets
.env
.env.local
.env.*.local
`

	if err := os.WriteFile(dir+"/.gitignore", []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("Error creating .gitignore: %v", err)
	}

	return nil
}

func createGinshotJSON(dir string, config *InitConfig) error {
	ginshotJSONContent := `{
	"project_name": "a",
	"port": "8080",
	"cors": true,
	"dockerfile": true,
	"docker_compose": true,
	"gitignore": true
}`

	if err := os.WriteFile(dir+"/ginshot.json", []byte(ginshotJSONContent), 0644); err != nil {
		return fmt.Errorf("Error creating ginshot.json: %v", err)
	}

	return nil
}

func createBrewerTemplateJSON(dir string, config *InitConfig) error {
	dataJSONContent := `{
	"project_name": "` + config.Name + `",
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
			"Data": "PingRequestData",
			"Message": "string"
		},
		"PingRequestData": {
			"Message": "string"
		}
	},
	"responses": {
		"PingResponse": {
			"Data": "PingResponseData",
			"Message": "string"
		},
		"PingResponseData": {
			"Ping": "models.Ping"
		}
	}
}`

	if err := os.WriteFile(dir+"/brewer/template.json", []byte(dataJSONContent), 0644); err != nil {
		return fmt.Errorf("Error creating template.json: %v", err)
	}

	return nil
}
