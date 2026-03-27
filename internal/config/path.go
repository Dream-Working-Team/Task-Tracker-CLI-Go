package config

import (
	"os"
	"path/filepath"
)

// ObtenerDirectorioDatos devuelve la ruta absoluta de la carpeta data local
func GetDirectionData() (string, error) {
	namneFolder := "data"

	// os.MkdirAll se asegura de que exista. Como ya la creaste en tu editor
	// simplemente pasará de largo sin dar error
	if err := os.MkdirAll(namneFolder, 0755); err != nil {
		return "", err
	}

	return filepath.Abs(namneFolder)
}
