package config

import (
	"os"
	"path/filepath"
)

// ObtenerDirectorioDatos devuelve la ruta absoluta de la carpeta "data" local.
func GetDirectionData() (string, error) {
	// Apuntamos directamente a la carpeta "data" que tienes en la raíz
	namneFolder := "data"

	// os.MkdirAll se asegura de que exista. Como ya la creaste en tu editor,
	// simplemente pasará de largo sin dar error.
	if err := os.MkdirAll(namneFolder, 0755); err != nil {
		return "", err
	}

	// filepath.Abs convierte "data" en la ruta completa de tu computadora
	// (ej. /home/antorlok/Documentos/Pruebas/test CLI/data)
	// Esto previene bugs si ejecutas el binario desde otra subcarpeta.
	return filepath.Abs(namneFolder)
}
