package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mates182/ginshot/formatter"
	"github.com/mates182/ginshot/models"
	"github.com/mates182/ginshot/reader"
	"github.com/spf13/cobra"
)

// TemplateData structure for models, requests, and responses
type TemplateData struct {
	Models    map[string]map[string]interface{} `json:"models"`
	Requests  map[string]map[string]interface{} `json:"requests"`
	Responses map[string]map[string]interface{} `json:"responses"`
	Messages  map[string]map[string]interface{} `json:"messages"`
}

// brewCmd represents the brew command
var brewCmd = &cobra.Command{
	Use:   "brew [filename]",
	Short: "Generate scaffold files from templates or JSON",
	Long:  `This command generates scaffold files from templates or JSON. For example, use a JSON to create models, requests, and responses.`,
	Run: func(cmd *cobra.Command, args []string) {
		brewerDir := "./brewer"
		if !directoryExists(brewerDir) {
			fmt.Println("Error: 'brewer' directory does not exist.")
			return
		}

		if len(args) > 0 {
			processTemplateFile(args[0], brewerDir)
		} else {
			listTemplates(brewerDir)
		}
	},
}

// processTemplateFile handles template file processing
func processTemplateFile(fileName, brewerDir string) {
	templatePath := filepath.Join(brewerDir, fileName)

	if !fileExists(templatePath) {
		fmt.Printf("Error: Template file '%s' does not exist in the 'brewer' directory.\n", fileName)
		return
	}

	if strings.HasSuffix(fileName, ".json") {
		generateScaffoldFromJSON(templatePath)
	} else {
		fmt.Printf("Generating scaffold from '%s' template...\n", fileName)
		content, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Printf("Error reading template file: %v\n", err)
			return
		}
		fmt.Println(string(content))
	}
}

// listTemplates displays available templates in the 'brewer' directory
func listTemplates(brewerDir string) {
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

	if len(templateFiles) == 0 {
		fmt.Println("No templates found in the 'brewer' directory.")
		var response string
		fmt.Print("Do you want to generate a new template? (y/n): ")
		fmt.Scanln(&response)

		if strings.ToLower(response) == "y" {
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

// generateScaffoldFromJSON generates models, requests, and responses from the JSON file
func generateScaffoldFromJSON(jsonPath string) {
	content, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	var template TemplateData
	if err = json.Unmarshal(content, &template); err != nil {
		fmt.Printf("Error parsing JSON file: %v\n", err)
		return
	}
	config, err := reader.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading project configuration: %v\n", err)
		return
	}
	projectConfig = config

	// Generate scaffold files
	//generateScaffold(template.Models, "./models/", "model", generateModelGo)
	//generateScaffold(template.Requests, "./data/requests/", "request", generateRequestGo)
	//generateScaffold(template.Responses, "./data/responses/", "response", generateResponseGo)
	//generateScaffold(template.Messages, "./data/messages/", "message", generateResponseGo)
}

// generateScaffold processes models, requests, or responses
func generateScaffold(data map[string]map[string]interface{}, directory, fileType string, config *models.ProjectConfig, generator func(string, map[string]interface{}, string) string) {
	if len(data) == 0 {
		return
	}

	for name, fields := range data {
		fileName := fmt.Sprintf("%s%s.go", directory, formatter.ToLowerCase(name))
		fileContent := generator(name, fields, config.ProjectName)
		if err := os.WriteFile(fileName, []byte(fileContent), 0644); err != nil {
			fmt.Printf("Error writing %s.go: %v\n", fileType, err)
			return
		}
		fmt.Printf("Scaffold for %s '%s' generated successfully!\n", fileType, name)
	}
}

// generateModelGo generates the Go code for models from the JSON structure
func generateModelGo(modelName string, fields map[string]interface{}, projectName string) string {
	var modelFields strings.Builder

	for fieldName, fieldType := range fields {
		tags := fmt.Sprintf("`json:\"%s\" bson:\"%s\"`", formatter.ToLowerCase(fieldName), formatter.ToLowerCase(fieldName))
		modelFields.WriteString(fmt.Sprintf("\t%s %s %s\n", fieldName, fieldType, tags))
	}

	return fmt.Sprintf(`package models

type %s struct {
%s}
`, modelName, modelFields.String())
}

// generateRequestGo generates the Go code for requests from the JSON structure
func generateRequestGo(requestName string, fields map[string]interface{}, projectName string) string {
	return generateStruct("request", requestName, fields, projectName)
}

// generateResponseGo generates the Go code for responses from the JSON structure
func generateResponseGo(responseName string, fields map[string]interface{}, projectName string) string {
	return generateStruct("response", responseName, fields, projectName)
}

// generateResponseGo generates the Go code for responses from the JSON structure
func generateMessageGo(responseName string, fields map[string]interface{}, projectName string) string {
	return generateStruct("message", responseName, fields, projectName)
}

// generateStruct generates struct code for requests and responses
func generateStruct(structType, structName string, fields map[string]interface{}, projectName string) string {
	var structFields strings.Builder
	needsImport := false

	for fieldName, fieldType := range fields {
		if strings.HasPrefix(fmt.Sprintf("%v", fieldType), "models.") || strings.HasPrefix(fmt.Sprintf("%v", fieldType), "[]models.") {
			needsImport = true
		}
		structFields.WriteString(fmt.Sprintf("\t%s %v `json:\"%s\"`\n", fieldName, fieldType, formatter.ToLowerCase(fieldName)))
	}

	importStatement := ""
	if needsImport {
		importStatement = fmt.Sprintf("import \"%s/internal/data/models\"\n", projectName)
	}

	return fmt.Sprintf(`package %s

%s
type %s struct {
%s}
`, structType, importStatement, structName, structFields.String())
}

// directoryExists checks if a directory exists
func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func init() {
	rootCmd.AddCommand(brewCmd)
}
