package storage

import (
	"encoding/json"
	"errors"
	"os"
)

func handleReadError[T any](err error) ([]T, error) {
	if errors.Is(err, os.ErrNotExist) {
		return []T{}, nil
	}
	return nil, err
}

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

func SaveFile[T any](ruta string, items []T) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ruta, data, 0644)
}

// Interfaces específicas exportadas
