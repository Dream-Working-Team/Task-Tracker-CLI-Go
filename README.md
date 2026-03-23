# Task Tracker CLI (Go)

Un gestor de tareas de interfaz de línea de comandos (CLI) simple, rápido y ligero construido completamente en Go. Este proyecto permite realizar un seguimiento de tus tareas diarias, gestionando sus estados (por hacer, en progreso, completadas) mediante almacenamiento persistente en un archivo JSON local.

## 🚀 Características

* **Cero Dependencias:** Construido exclusivamente con la biblioteca estándar de Go (`os`, `encoding/json`, `fmt`, `time`).
* **Almacenamiento Local:** Las tareas se guardan automáticamente en un archivo `tasks.json` en el directorio actual. El archivo se crea automáticamente si no existe.
* **Gestión Completa (CRUD):** Permite añadir, actualizar y eliminar tareas fácilmente.
* **Control de Estados:** Transiciones fluidas entre estados (`todo`, `in-progress`, `done`).
* **Filtros de Búsqueda:** Listado de tareas global o filtrado por su estado actual.
* **Manejo de Errores Robusto:** Respuestas claras ante comandos inválidos, IDs inexistentes o problemas de permisos en el sistema de archivos.

## 🛠️ Instalación y Compilación

Asegúrate de tener [Go instalado](https://go.dev/doc/install) en tu sistema.

1. Clona este repositorio:
   ```bash
   git clone [https://github.com/TU_ORGANIZACION/task-tracker-go.git](https://github.com/TU_ORGANIZACION/task-tracker-go.git)
   cd task-tracker-go
