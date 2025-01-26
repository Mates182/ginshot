package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// La estructura del JSON de entrada
type TemplateData struct {
	Models   map[string]map[string]interface{} `json:"models"`
	Requests map[string]map[string]interface{} `json:"requests"`
}

// brewCmd representa el comando brew
var brewCmd = &cobra.Command{
	Use:   "brew [filename]",
	Short: "Generate scaffold files from templates or JSON",
	Long:  `This command generates scaffold files from templates or JSON. For example, use a JSON to create models and requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Directorio brewer
		brewerDir := "./brewer"
		if _, err := os.Stat(brewerDir); os.IsNotExist(err) {
			fmt.Println("Error: 'brewer' directory does not exist.")
			return
		}

		// Si se proporciona un nombre de archivo, generar el scaffold a partir de él
		if len(args) > 0 {
			fileName := args[0]
			templatePath := filepath.Join(brewerDir, fileName)

			// Verificar si el archivo existe
			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				fmt.Printf("Error: Template file '%s' does not exist in the 'brewer' directory.\n", fileName)
				return
			}

			// Procesar archivos JSON para generar modelos y solicitudes
			if strings.HasSuffix(fileName, ".json") {
				generateScaffoldFromJSON(templatePath)
			} else {
				// Manejar otros archivos de plantilla (por ejemplo, .txt u otros formatos)
				fmt.Printf("Generating scaffold from '%s' template...\n", fileName)
				content, err := os.ReadFile(templatePath)
				if err != nil {
					fmt.Printf("Error reading template file: %v\n", err)
					return
				}

				// Por ahora, solo imprimir el contenido (más tarde puedes procesarlo)
				fmt.Println(string(content))
			}
		} else {
			// Listar las plantillas en el directorio 'brewer'
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
					// Lógica para generar una nueva plantilla
					fmt.Println("Generating a new template...")
					// Añadir lógica para generar una plantilla aquí
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

// generateScaffoldFromJSON genera los modelos y requests a partir del JSON
func generateScaffoldFromJSON(jsonPath string) {
	// Leer el archivo JSON
	content, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	// Parsear el contenido del JSON
	var template TemplateData
	err = json.Unmarshal(content, &template)
	if err != nil {
		fmt.Printf("Error parsing JSON file: %v\n", err)
		return
	}

	// Si existen modelos, generarlos
	if len(template.Models) > 0 {
		for modelName, modelFields := range template.Models {
			modelFileName := fmt.Sprintf("./models/%s.go", modelName)
			modelFileContent := generateModelGo(modelName, modelFields)
			err = os.WriteFile(modelFileName, []byte(modelFileContent), 0644)
			if err != nil {
				fmt.Printf("Error writing model.go: %v\n", err)
				return
			}

			// Mensaje de éxito para el modelo
			fmt.Printf("Scaffold for model '%s' generated successfully!\n", modelName)
		}
	}

	// Si existen requests, generarlos
	if len(template.Requests) > 0 {
		for requestName, requestFields := range template.Requests {
			requestFileName := fmt.Sprintf("./data/requests/%s.go", requestName)
			requestFileContent := generateRequestGo(requestName, requestFields)
			err = os.WriteFile(requestFileName, []byte(requestFileContent), 0644)
			if err != nil {
				fmt.Printf("Error writing request.go: %v\n", err)
				return
			}

			// Mensaje de éxito para el request
			fmt.Printf("Scaffold for request '%s' generated successfully!\n", requestName)
		}
	}
}

// generateModelGo crea el código Go para los modelos a partir de la estructura JSON
func generateModelGo(modelName string, fields interface{}) string {
	var modelFields string

	switch v := fields.(type) {
	case map[string]interface{}:
		// Si es un objeto anidado, procesar los campos recursivamente
		for fieldName, fieldType := range v {
			// Verificar si el campo es un struct (objeto anidado)
			if nestedField, ok := fieldType.(map[string]interface{}); ok {
				// Crear recursivamente un struct para el campo anidado
				modelFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldName, fieldName)
				// Generar el modelo anidado
				modelFields += generateModelGo(fieldName, nestedField)
			} else {
				// Campo simple (string, int, etc.)
				modelFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, fieldName)
			}
		}
	}

	return fmt.Sprintf(`package models

type %s struct {
%s}
`, modelName, modelFields)
}

// generateRequestGo crea el código Go para los requests a partir de la estructura JSON
func generateRequestGo(requestName string, fields interface{}) string {
	var requestFields string

	switch v := fields.(type) {
	case map[string]interface{}:
		// Si es un objeto anidado, procesar los campos recursivamente
		for fieldName, fieldType := range v {
			// Verificar si el campo es un struct (objeto anidado)
			if nestedField, ok := fieldType.(map[string]interface{}); ok {
				// Crear recursivamente un struct para el campo anidado
				requestFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldName, fieldName)
				// Generar el request anidado
				requestFields += generateRequestGo(fieldName, nestedField)
			} else {
				// Campo simple (string, int, etc.)
				requestFields += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, fieldName)
			}
		}
	}

	return fmt.Sprintf(`package request

type %s struct {
%s}
`, requestName, requestFields)
}

func init() {
	rootCmd.AddCommand(brewCmd)
}
