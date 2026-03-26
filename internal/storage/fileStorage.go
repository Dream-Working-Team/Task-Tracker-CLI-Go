package storage

// Realizamos la importaciones correspondientes
// para manejar jsons, errores, funciones del sistema, y el modelo previamente creado
import (
	"encoding/json"
	"errors"
	"os"
	"task-cli/internal/model"
)

// Creamos la estructura del alamcenamiento con su ruta
type storage struct {
	filename string
}

// Configuramos la instancia
// Utilizamos *storage para devolver un puntero a la instancia
// Se usa &storage para que se tome la direccion en memoria
func newStorage(root string) *storage {
	return &storage{
		filename: root,
	}
}

// Para manejar los errores de manera modularizada
// Esta funcion recibe una lista del modelo tareas y el error
// El error.Is compara los errores, si lo cumple devuelve la lista vacia y nil sin error, para la primera
// vez que se abra el programa, sino devuelve nil y el error
func handleReadError(err error) ([]model.Task, error) {
	if errors.Is(err, os.ErrNotExist) {

		return []model.Task{}, nil

	}

	return nil, err
}

// Lee los archivos para convertirlo en objetos
// Tiene acceso a la instancia actual de storage para el filename
// Devolvemos la lista vacia, llena o con errores si ocurrio alguno
func (a *storage) readTask() ([]model.Task, error) {

	data, err := os.ReadFile(a.filename)

	if err != nil {
		return handleReadError(err)

	}

	if len(data) == 0 {
		return []model.Task{}, nil
	}

	var tasks []model.Task

	err = json.Unmarshal(data, &tasks)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Guarda los datos en el json con indentacion
// De lo contrario retorna un error
// Dato: el 0644, es para Linux los permisos de lectura y escritura
func (a *storage) saveTask(Tasks []model.Task) error {

	data, err := json.MarshalIndent(Tasks, "", " ")

	if err != nil {
		return err
	}

	return os.WriteFile(a.filename, data, 0644)

}
