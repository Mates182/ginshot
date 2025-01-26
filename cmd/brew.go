package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Field struct to hold the field name and its type
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Model struct to hold the model definition
type Model struct {
	Fields map[string]interface{} `json:"fields"`
}

// The structure of the JSON input
type Models struct {
	Models map[string]map[string]interface{} `json:"models"`
}

// brewCmd represents the brew command
var brewCmd = &cobra.Command{
	Use:   "brew [filename]",
	Short: "Generate scaffold files from templates or JSON",
	Long:  `This command generates scaffold files from templates or JSON. For example, use a JSON to create models.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Brewer directory
		brewerDir := "./brewer"
		if _, err := os.Stat(brewerDir); os.IsNotExist(err) {
			fmt.Println("Error: 'brewer' directory does not exist.")
			return
		}

		// If a filename is provided, generate scaffold from it
		if len(args) > 0 {
			fileName := args[0]
			templatePath := filepath.Join(brewerDir, fileName)

			// Check if the file exists
			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				fmt.Printf("Error: Template file '%s' does not exist in the 'brewer' directory.\n", fileName)
				return
			}

			// Process JSON files for model scaffolding
			if strings.HasSuffix(fileName, ".json") {
				generateModelsFromJSON(templatePath)
			} else {
				// Handle other template files (e.g., .txt or other formats)
				fmt.Printf("Generating scaffold from '%s' template...\n", fileName)
				content, err := ioutil.ReadFile(templatePath)
				if err != nil {
					fmt.Printf("Error reading template file: %v\n", err)
					return
				}

				// For now, print the content (later you can process it)
				fmt.Println(string(content))
			}
		} else {
			// List templates in the 'brewer' directory
			files, err := ioutil.ReadDir(brewerDir)
			if err != nil {
				fmt.Println("Error reading 'brewer' directory:", err)
				return
			}

			var templateFiles []string
			for _, file := range files {
				if !file.IsDir() {
					templateFiles = append(templateFiles, file.Name())
				}
			}

			if len(templateFiles) == 0 {
				fmt.Println("No templates found in the 'brewer' directory.")
				var response string
				fmt.Print("Do you want to generate a new template? (y/n): ")
				fmt.Scanln(&response)

				if strings.ToLower(response) == "y" {
					// Logic for generating a new template
					fmt.Println("Generating a new template...")
					// Add logic to generate a template file here
				} else {
					fmt.Println("Exiting. No template will be created.")
				}
			} else {
				fmt.Println("Available templates in 'brewer' directory:")
				for _, file := range templateFiles {
					fmt.Println(" -", file)
				}
			}
		}
	},
}

// generateModelsFromJSON generates model files based on the JSON structure
func generateModelsFromJSON(jsonPath string) {
	// Read the JSON file
	content, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	// Parse the JSON content into the Models struct
	var models Models
	err = json.Unmarshal(content, &models)
	if err != nil {
		fmt.Printf("Error parsing JSON file: %v\n", err)
		return
	}

	// Generate models for each entry in the models
	for modelName, modelFields := range models.Models {
		modelFileName := fmt.Sprintf("./models/%s.go", modelName)
		modelFileContent := generateModelGo(modelName, modelFields)
		err = ioutil.WriteFile(modelFileName, []byte(modelFileContent), 0644)
		if err != nil {
			fmt.Printf("Error writing model.go: %v\n", err)
			return
		}

		// Success message
		fmt.Printf("Scaffold for model '%s' generated successfully!\n", modelName)
	}
}

// generateModelGo creates the Go code for the model based on the JSON structure
func generateModelGo(modelName string, fields interface{}) string {
	var modelFields string

	switch v := fields.(type) {
	case map[string]interface{}:
		// If it's a nested object, recursively process the fields
		for fieldName, fieldType := range v {
			// Check if the field is an object itself (nested structure)
			if nestedField, ok := fieldType.(map[string]interface{}); ok {
				// Recursively create a struct for the nested field
				modelFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldName, fieldName)
				// Generate the nested model
				modelFields += generateModelGo(fieldName, nestedField)
			} else {
				modelFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, fieldName)
			}
		}
	}

	return fmt.Sprintf(`package models

type %s struct {
%s}
`, modelName, modelFields)
}

func init() {
	rootCmd.AddCommand(brewCmd)
}
