package reader

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mates182/ginshot/models"
)

func LoadConfig() (*models.ProjectConfig, error) {
	filePath := "ginshot.json"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %v", err)
	}
	defer file.Close()

	var config models.ProjectConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error parsing the JSON file: %v", err)
	}

	return &config, nil
}
