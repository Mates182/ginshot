package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/mates182/ginshot/models"
	"github.com/mates182/ginshot/reader"
	templates "github.com/mates182/ginshot/templ"
	"github.com/mates182/ginshot/writer"
	"github.com/spf13/cobra"
)

// pourCmd represents the pour command
var pourCmd = &cobra.Command{
	Use:   "pour",
	Short: "Generate files for the selected database (Redis or MongoDB) and a custom name",
	Long: `This command generates the necessary files for your project 
to integrate with the specified database, such as Redis or MongoDB, 
and applies a custom name to the generated files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the database flag
		db, _ := cmd.Flags().GetString("db")
		// Get the name flag
		name, _ := cmd.Flags().GetString("name")

		// Check if a database type was provided
		if db != "" {
			// If name is provided, use it; otherwise, prompt for a name
			if name != "" {
				// Generate files based on the database and name
				createDBFiles(db, name)
			} else {
				// Ask the user for the name
				promptForName(db)
			}
		} else {
			// Ask the user for the database type if not provided
			promptForDB()
		}
	},
}

func init() {
	// Register the pour command
	rootCmd.AddCommand(pourCmd)

	// Add flags for specifying the database and name
	pourCmd.Flags().String("db", "", "Database to use (redis, mongo): ")
	pourCmd.Flags().String("name", "", "Database name: ")
}

// ReadProjectName lee el archivo ginshot.json y devuelve el nombre del proyecto
func ReadProjectName() (string, error) {
	// Abrir el archivo ginshot.json
	file, err := os.Open("ginshot.json")
	if err != nil {
		return "", fmt.Errorf("error opening ginshot.json: %v", err)
	}
	defer file.Close()

	// Decodificar el archivo JSON
	var config models.ProjectConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return "", fmt.Errorf("error decoding ginshot.json: %v", err)
	}

	// Devolver el nombre del proyecto
	return config.ProjectName, nil
}

// createDBFiles generates the necessary files based on the chosen database and name
func createDBFiles(db, name string) {
	config, err := reader.LoadConfig()
	if err != nil {
		fmt.Println("Error reading project config:", err)
		return
	}
	config.Database.Name = name
	config.Database.Type = db
	config.Database.Table = db

	writer.SaveConfig(config)
	switch db {
	case "redis":
		// Create Redis files with the specified name
		fmt.Printf("Generating Redis files with the name '%s'...\n", name)
		// Add your file creation logic here, using the name (e.g., create a Redis client config with the name)
		// os.Create(name + "_redis_config.go")
	case "mongo":
		// Create MongoDB files with the specified name
		fmt.Printf("Generating MongoDB files with the name '%s'...\n", name)
		// Add your file creation logic here, using the name (e.g., create a MongoDB client config with the name)
		// os.Create(name + "_mongo_config.go")
	default:
		fmt.Println("Unknown database type. Please choose either redis or mongo.")
	}
	// Generate dbcontext files
	generateDBContextFiles(config, db, name)
	// Pass the name to the function
	generateSecretsFile(db, name) // Pass the name to the function

	generateEnvFile(db, name)

	addDatabaseToDockerCompose(db, name)

	addDBConfig(db, name)

	// Clean up Go modules
	if err := runGoModTidy("."); err != nil {
		fmt.Println(err)
		return
	}
}

// promptForDB asks the user which database to use (redis or mongo) if none was provided
func promptForDB() {
	// Ask the user for the database type
	var db string
	fmt.Println("No database specified.")
	fmt.Print("Please choose a database (redis or mongo): ")
	_, err := fmt.Scanln(&db)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Generate files based on the user's choice
	promptForName(db)
}

// promptForName asks the user for a custom name if it was not provided
func promptForName(db string) {
	// Ask the user for a custom name
	var name string
	fmt.Println("No name specified.")
	fmt.Print("Please provide a custom name for the generated files: ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Generate files based on the chosen database and name
	createDBFiles(db, name)
}

func generateDBContextFiles(config *models.ProjectConfig, db, name string) {
	// Define the base directory for dbcontext
	baseDir := fmt.Sprintf("dbcontext/%s", name)

	// Create the directory for dbcontext
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Create dbcontext.go for the selected database (mongo or redis)
	dbContextTemplate := templates.GetDBContextTemplate(config, db, name)

	// Generate the dbcontext.go file content
	tmpl, err := template.New("dbcontext").Parse(dbContextTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Create the dbcontext.go file
	dbContextFilePath := fmt.Sprintf("%s/dbcontext.go", baseDir)
	file, err := os.Create(dbContextFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Execute the template and write it to the file
	err = tmpl.Execute(file, nil)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Printf("Successfully generated %s/dbcontext.go for %s\n", baseDir, db)
}

func generateSecretsFile(db string, name string) {
	// Define the secrets directory
	secretsDir := "secrets"
	err := os.MkdirAll(secretsDir, 0755)
	if err != nil {
		fmt.Println("Error creating secrets directory:", err)
		return
	}

	// Determine the template based on the database
	var secretsTemplate string
	switch db {
	case "mongo":
		secretsTemplate = `
package secrets

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Get` + name + `DBURI() string {
	err := godotenv.Load()
	if err != nil {
		// TODO: remove on production
		fmt.Println("Error loading .env file, ignore if is on docker")
	}
	dbUri := os.Getenv("` + name + `_URI")
	return dbUri
}
`
	case "redis":
		secretsTemplate = `
package secrets

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Get` + name + `DBURI() string {
	err := godotenv.Load()
	if err != nil {
		// TODO: remove on production
		fmt.Println("Error loading .env file, ignore if is on docker")
	}
	dbUri := os.Getenv("` + name + `_URI")
	return dbUri
}

func Get` + name + `DBPassword() string {
	err := godotenv.Load()
	if err != nil {
		// TODO: remove on production
		fmt.Println("Error loading .env file, ignore if is on docker")
	}
	dbPassword := os.Getenv("` + name + `_PASSWORD")
	return dbPassword
}
`
	default:
		fmt.Println("Unsupported database type.")
		return
	}

	// Generate the env.go file content
	tmpl, err := template.New("secrets").Parse(secretsTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Create the env.go file
	secretsFilePath := fmt.Sprintf("%s/%s-env.go", secretsDir, name)
	file, err := os.Create(secretsFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Execute the template and write it to the file
	err = tmpl.Execute(file, nil)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Printf("Successfully generated %s/env.go\n", secretsDir)
}

func generateEnvFile(db, name string) {
	// Define the template based on the database type
	var envTemplate string
	switch db {
	case "mongo":
		envTemplate = `# MongoDB Connection URI for ` + name + `
` + name + `_URI=mongodb://localhost:27017/` + name

	case "redis":
		envTemplate = `# Redis Connection Details for ` + name + `
` + name + `_URI=localhost:6379
` + name + `_PASSWORD=`

	default:
		fmt.Println("Unsupported database type.")
		return
	}

	// Create the .env file
	file, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating .env file:", err)
		return
	}
	defer file.Close()

	// Add a newline before appending if file already exists and doesn't end with one
	info, err := file.Stat()
	if err == nil && info.Size() > 0 {
		file.WriteString("\n")
	}

	// Write the environment variables
	_, err = file.WriteString(envTemplate + "\n")
	if err != nil {
		fmt.Println("Error writing to .env file:", err)
		return
	}

	fmt.Println("Successfully updated .env file with " + name + " configuration")
}

// addDatabaseToDockerCompose adds MongoDB or Redis service to the docker-compose.yml based on the provided database type
func addDatabaseToDockerCompose(dbType, name string) {
	// Read the docker-compose.yml file
	filePath := "./docker-compose.yml"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading docker-compose.yml:", err)
		return
	}

	// Convert the file content to a string
	content := string(fileContent)

	// Define the MongoDB and Redis service templates
	mongoService := `
  ` + name + `-db:
    build:
      context: ./config/` + name + `
      dockerfile: Dockerfile
    container_name: ` + name + `-db
    environment:
      MONGO_INITDB_ROOT_USERNAME: user
      MONGO_INITDB_ROOT_PASSWORD: password
    ports:
      - "27017:27017"
    networks:
      - app-network
`

	redisService := `
  ` + name + `-db:
    image: redis:7.0-alpine
    container_name: ` + name + `-db
    ports:
      - "6379:6379"
    networks:
      - app-network
    environment:
      - REDIS_PASSWORD=password
`

	// Define environment variables for the application
	mongoEnv := `
      - ` + name + `_URI=mongodb://user:password@` + name + `-db:27017`

	redisEnv := `
      - ` + name + `_URI=` + name + `-db:6379
      - ` + name + `_PASSWORD=password`

	// Find the position of the "services:" line and add the service just below it
	servicesPosition := strings.Index(content, "services:")
	if servicesPosition == -1 {
		fmt.Println("Error: 'services:' section not found in docker-compose.yml.")
		return
	}

	// Add the corresponding service after the "services:" section
	if dbType == "mongo" {
		if !strings.Contains(content, name+"-db") {
			content = content[:servicesPosition+len("services:")] + mongoService + content[servicesPosition+len("services:"):]
			content = strings.Replace(content, "- GIN_MODE=release", "- GIN_MODE=release"+mongoEnv, 1)
			fmt.Println("MongoDB service added to docker-compose.yml.")
		} else {
			fmt.Println("MongoDB service already exists in docker-compose.yml.")
		}
	} else if dbType == "redis" {
		if !strings.Contains(content, name+"-db") {
			content = content[:servicesPosition+len("services:")] + redisService + content[servicesPosition+len("services:"):]
			content = strings.Replace(content, "- GIN_MODE=release", "- GIN_MODE=release"+redisEnv, 1)
			fmt.Println("Redis service added to docker-compose.yml.")
		} else {
			fmt.Println("Redis service already exists in docker-compose.yml.")
		}
	} else {
		fmt.Println("Unsupported database type. Please choose either 'mongo' or 'redis'.")
		return
	}

	// Add depends_on after restart: unless-stopped
	dependsOnLine := "\n    depends_on:\n      - " + name + "-db"
	content = strings.Replace(content, "restart: unless-stopped", "restart: unless-stopped"+dependsOnLine, 1)

	// Write the updated content back to the docker-compose.yml file
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing to docker-compose.yml:", err)
		return
	}

	// Indicate success
	fmt.Println("docker-compose.yml updated successfully.")
}

func addDBConfig(db string, name string) {
	if db == "mongo" {
		dirPath := "config/" + name
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Println("Error creating config directory:", err)
			return
		}
		var collection_name string
		fmt.Print("Enter a collection name: ")
		fmt.Scanln(&collection_name)

		initDBScript := `#!/bin/bash
mongoimport --host localhost --db ` + name + ` --collection ` + collection_name + ` --file /data/backup.json --jsonArray
echo "Database initialized successfully"`

		backupData := `[{
	"Message": "Pong",
	"Greeting": {
		"Hello": "Hello World"
	}
}]`
		dockerfile := `FROM mongo:latest

WORKDIR /data

COPY backup.json /data/backup.json

COPY init-db.sh /docker-entrypoint-initdb.d/init-db.sh

RUN chmod +x /docker-entrypoint-initdb.d/init-db.sh

EXPOSE 27017`

		// init-db.sh
		err = os.WriteFile(dirPath+"/init-db.sh", []byte(initDBScript), 0755)
		if err != nil {
			fmt.Println("Error creating init-db.sh:", err)
			return
		}

		// backup.json
		err = os.WriteFile(dirPath+"/backup.json", []byte(backupData), 0644)
		if err != nil {
			fmt.Println("Error creating backup.json:", err)
			return
		}
		// dockerfile
		err = os.WriteFile(dirPath+"/Dockerfile", []byte(dockerfile), 0644)
		if err != nil {
			fmt.Println("Error creating Dockerfile", err)
			return
		}

		fmt.Println("MongoDB configuration files created successfully.")
	} else {
		fmt.Println("No additional files required for the selected database.")
	}
}
