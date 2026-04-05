package config

import (
	"os"
	"path/filepath"
)

// GetDataDirectory returns the absolute path to the local data folder
func GetDataDirectory() (string, error) {
	dataFolder := "data"

	// os.MkdirAll ensures the folder exists
	// If it already exists, it continues without returning an error
	if err := os.MkdirAll(dataFolder, 0755); err != nil {
		return "", err
	}

	return filepath.Abs(dataFolder)
}
