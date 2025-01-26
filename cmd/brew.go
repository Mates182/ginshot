package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// TemplateData structure for models, requests, and responses
type TemplateData struct {
	Models    map[string]map[string]interface{} `json:"models"`
	Requests  map[string]map[string]interface{} `json:"requests"`
	Responses map[string]map[string]interface{} `json:"responses"`
}

// brewCmd represents the brew command
var brewCmd = &cobra.Command{
	Use:   "brew [filename]",
	Short: "Generate scaffold files from templates or JSON",
	Long:  `This command generates scaffold files from templates or JSON. For example, use a JSON to create models, requests, and responses.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Brewer directory where templates are stored
		brewerDir := "./brewer"
		if _, err := os.Stat(brewerDir); os.IsNotExist(err) {
			fmt.Println("Error: 'brewer' directory does not exist.")
			return
		}

		// If a filename is provided, generate the scaffold from it
		if len(args) > 0 {
			fileName := args[0]
			templatePath := filepath.Join(brewerDir, fileName)

			// Check if the file exists
			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				fmt.Printf("Error: Template file '%s' does not exist in the 'brewer' directory.\n", fileName)
				return
			}

			// If it's a JSON file, process it to generate scaffold
			if strings.HasSuffix(fileName, ".json") {
				generateScaffoldFromJSON(templatePath)
			} else {
				// Handle other template file types (e.g., .txt, etc.)
				fmt.Printf("Generating scaffold from '%s' template...\n", fileName)
				content, err := os.ReadFile(templatePath)
				if err != nil {
					fmt.Printf("Error reading template file: %v\n", err)
					return
				}

				// Currently just print the template content (you can extend this logic)
				fmt.Println(string(content))
			}
		} else {
			// List templates in the 'brewer' directory if no file is provided
			files, err := os.ReadDir(brewerDir)
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

			// If no templates found, ask if the user wants to create one
			if len(templateFiles) == 0 {
				fmt.Println("No templates found in the 'brewer' directory.")
				var response string
				fmt.Print("Do you want to generate a new template? (y/n): ")
				fmt.Scanln(&response)

				if strings.ToLower(response) == "y" {
					// Logic for generating a new template (this can be extended)
					fmt.Println("Generating a new template...")
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

// generateScaffoldFromJSON generates models, requests, and responses from the JSON file
func generateScaffoldFromJSON(jsonPath string) {
	// Read the JSON file
	content, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	// Parse the JSON content
	var template TemplateData
	err = json.Unmarshal(content, &template)
	if err != nil {
		fmt.Printf("Error parsing JSON file: %v\n", err)
		return
	}

	// Generate models if they exist
	if len(template.Models) > 0 {
		for modelName, modelFields := range template.Models {
			modelFileName := fmt.Sprintf("./models/%s.go", modelName)
			modelFileContent := generateModelGo(modelName, modelFields)
			err = os.WriteFile(modelFileName, []byte(modelFileContent), 0644)
			if err != nil {
				fmt.Printf("Error writing model.go: %v\n", err)
				return
			}

			// Success message for model
			fmt.Printf("Scaffold for model '%s' generated successfully!\n", modelName)
		}
	}

	// Generate requests if they exist
	if len(template.Requests) > 0 {
		for requestName, requestFields := range template.Requests {
			requestFileName := fmt.Sprintf("./data/requests/%s.go", requestName)
			requestFileContent := generateRequestGo(requestName, requestFields)
			err = os.WriteFile(requestFileName, []byte(requestFileContent), 0644)
			if err != nil {
				fmt.Printf("Error writing request.go: %v\n", err)
				return
			}

			// Success message for request
			fmt.Printf("Scaffold for request '%s' generated successfully!\n", requestName)
		}
	}

	// Generate responses if they exist
	if len(template.Responses) > 0 {
		for responseName, responseFields := range template.Responses {
			responseFileName := fmt.Sprintf("./data/responses/%s.go", responseName)
			responseFileContent := generateResponseGo(responseName, responseFields)
			err = os.WriteFile(responseFileName, []byte(responseFileContent), 0644)
			if err != nil {
				fmt.Printf("Error writing response.go: %v\n", err)
				return
			}

			// Success message for response
			fmt.Printf("Scaffold for response '%s' generated successfully!\n", responseName)
		}
	}
}

// generateModelGo generates the Go code for models from the JSON structure
func generateModelGo(modelName string, fields interface{}) string {
	var modelFields string

	switch v := fields.(type) {
	case map[string]interface{}:
		// If it's a nested object, process the fields recursively
		for fieldName, fieldType := range v {
			// Check if the field is a nested struct (object)
			if nestedField, ok := fieldType.(map[string]interface{}); ok {
				// Create a struct for the nested field
				modelFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldName, fieldName)
				// Generate nested model
				modelFields += generateModelGo(fieldName, nestedField)
			} else {
				// Simple field (string, int, etc.)
				modelFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, fieldName)
			}
		}
	}

	return fmt.Sprintf(`package models

type %s struct {
%s}
`, modelName, modelFields)
}

// generateRequestGo generates the Go code for requests from the JSON structure
func generateRequestGo(requestName string, fields interface{}) string {
	var requestFields string

	switch v := fields.(type) {
	case map[string]interface{}:
		// If it's a nested object, process the fields recursively
		for fieldName, fieldType := range v {
			// Check if the field is a nested struct (object)
			if nestedField, ok := fieldType.(map[string]interface{}); ok {
				// Create a struct for the nested field
				requestFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldName, fieldName)
				// Generate nested request
				requestFields += generateRequestGo(fieldName, nestedField)
			} else {
				// Simple field (string, int, etc.)
				requestFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, fieldName)
			}
		}
	}

	return fmt.Sprintf(`package request

type %s struct {
%s}
`, requestName, requestFields)
}

// generateResponseGo generates the Go code for responses from the JSON structure
func generateResponseGo(responseName string, fields interface{}) string {
	var responseFields string

	switch v := fields.(type) {
	case map[string]interface{}:
		// If it's a nested object, process the fields recursively
		for fieldName, fieldType := range v {
			// Check if the field is a nested struct (object)
			if nestedField, ok := fieldType.(map[string]interface{}); ok {
				// Create a struct for the nested field
				responseFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldName, fieldName)
				// Generate nested response
				responseFields += generateResponseGo(fieldName, nestedField)
			} else {
				// Simple field (string, int, etc.)
				responseFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, fieldName)
			}
		}
	}

	return fmt.Sprintf(`package response

type %s struct {
%s}
`, responseName, responseFields)
}

func init() {
	rootCmd.AddCommand(brewCmd)
}
