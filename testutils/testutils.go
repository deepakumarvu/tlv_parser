package testutils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func InterfaceToStruct(data interface{}, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling map to JSON: %w", err)
	}

	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON to struct: %w", err)
	}

	return nil
}

func GetFilePathInWD(file string) string {
	currentFilePath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current file path:", err)
		return currentFilePath
	}

	// Construct the file path
	return filepath.Join(currentFilePath, file)

}

func ReadJsonToStruct(filePath string, output interface{}) error {

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return err
	}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(content, output)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}
	return nil
}
