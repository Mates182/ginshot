package writer

import (
	"fmt"
	"os"
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
