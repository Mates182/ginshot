package writer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mates182/ginshot/models"
)

func WriteFile(filePath, content string) error {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error creating %s: %v", filePath, err)
	}
	return nil
}

func CreateDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %v", dirPath, err)
	}
	return nil
}

func SaveConfig(dir string, config *models.ProjectConfig) error {
	filePath := dir + "/ginshot.json"

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating the file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding the JSON file: %v", err)
	}

	return nil
}
