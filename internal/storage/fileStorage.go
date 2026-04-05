package storage

import (
	"encoding/json"
	"errors"
	"os"
)

// handleReadError converts a missing file error into an empty list and forwards other errors
func handleReadError[T any](err error) ([]T, error) {
	if errors.Is(err, os.ErrNotExist) {
		return []T{}, nil
	}
	return nil, err
}

// ReadFile reads a JSON file and converts it into a typed list
func ReadFile[T any](route string) ([]T, error) {
	data, err := os.ReadFile(route)
	if err != nil {
		return handleReadError[T](err)
	}
	if len(data) == 0 {
		return []T{}, nil
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// SaveFile writes a typed list to disk as indented JSON
func SaveFile[T any](route string, items []T) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(route, data, 0644)
}
