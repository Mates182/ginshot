package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

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
	pourCmd.Flags().String("db", "", "Database to use (redis, mongo)")
	pourCmd.Flags().String("name", "", "Custom name for the generated files")
}

// ReadProjectName lee el archivo ginshot.json y devuelve el nombre del proyecto
func ReadProjectName() (string, error) {
	// Abrir el archivo ginshot.json
	file, err := os.Open("ginshot.json")
	if err != nil {
		return "", fmt.Errorf("Error opening ginshot.json: %v", err)
	}
	defer file.Close()

	// Decodificar el archivo JSON
	var config ProjectConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return "", fmt.Errorf("Error decoding ginshot.json: %v", err)
	}

	// Devolver el nombre del proyecto
	return config.ProjectName, nil
}

// createDBFiles generates the necessary files based on the chosen database and name
func createDBFiles(db, name string) {
	projectName, err := ReadProjectName()
	if err != nil {
		fmt.Println("Error reading project name:", err)
		return
	}
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
	generateDBContextFiles(projectName, db, name)
	// Pass the name to the function
	generateSecretsFile(db, name) // Pass the name to the function

	generateEnvFile(db, name)

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

func generateDBContextFiles(projectName, db, name string) {
	// Define the base directory for dbcontext
	baseDir := fmt.Sprintf("dbcontext/%s", name)

	// Create the directory for dbcontext
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Create dbcontext.go for the selected database (mongo or redis)
	var dbContextTemplate string

	switch db {
	case "mongo":
		dbContextTemplate = `
package dbcontext

import (
	"context"
	"` + projectName + `/secrets"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var clientInstance *mongo.Client
var clientOnce sync.Once

func GetDBClient() *mongo.Client {
	clientOnce.Do(func() {
		endpoint := secrets.Get` + name + `DBURI()
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(endpoint))
		if err != nil {
			panic(err)
		}
		fmt.Println("Connected to ` + name + ` Database Server")

		err = client.Ping(context.Background(), readpref.Primary())
		if err != nil {
			panic(err)
		}
		fmt.Println("Pong")
		clientInstance = client
	})

	return clientInstance
}
		`
	case "redis":
		dbContextTemplate = `
package dbcontext

import (
	"context"
	"` + projectName + `/secrets"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func GetDBClient() *redis.Client {

	dbURI := secrets.Get` + name + `DBURI()
	dbPassword := secrets.Get` + name + `DBPassword()
	dbOptions := &redis.Options{
		Addr: dbURI,
		DB:   0,
	}
	if dbPassword != "" {
		dbOptions.Password = dbPassword
	}

	client := redis.NewClient(dbOptions)
	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to DB: %s\n", err.Error())
		return nil
	}
	fmt.Printf("Ping: %s\n", ping)
	return client
}
`
	default:
		fmt.Println("Unsupported database type.")
		return
	}

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
