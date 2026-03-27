package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"task-cli/internal/model"
	"task-cli/internal/storage"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"path/filepath"

	"task-cli/internal/config"
)

// getRouteSesion devuelve la ruta del archivo local donde se guarda la sesión.
func getRouteSesion() string {
	dir, err := config.GetDirectionData()
	if err != nil {
		return ".session" // Fallback de emergencia a la raíz local
	}
	return filepath.Join(dir, ".session")
}

// Valores más altos mejoran la seguridad, pero aumentan el tiempo de procesamiento
const passwordHashCost = bcrypt.DefaultCost

// AuthService coordina el registro de usuarios, el inicio de sesión y la persistencia de sesión
type AuthService struct {
	save *storage.Storage
}

// NewAuthService crea una instancia del servicio de autenticación.
func NewAuthService(s *storage.Storage) *AuthService {
	return &AuthService{save: s}
}

// Register registra un usuario nuevo si el nombre no existe.
func (s *AuthService) Register(username, password string) error {
	usersGeneral, err := s.save.ReadUsers()
	if err != nil {
		return err
	}

	newID := 1
	if len(usersGeneral) > 0 {
		newID = usersGeneral[len(usersGeneral)-1].ID + 1
	}

	for _, u := range usersGeneral {
		if u.Username == username {
			return errors.New("el usuario ya existe")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usersGeneral = append(usersGeneral, model.User{
		ID:           newID, // Asignamos el ID
		Username:     username,
		PasswordHash: string(hash),
	})
	return s.save.SaveUsers(usersGeneral)
}

// Login valida credenciales y guarda la sesión del usuario autenticado.
func (s *AuthService) Login(username, password string) error {
	usersGeneral, err := s.save.ReadUsers()
	if err != nil {
		return err
	}

	for _, u := range usersGeneral {
		if u.Username != username {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
			return errors.New("credenciales inválidas")
		}
		idStr := strconv.Itoa(u.ID)
		return os.WriteFile(getRouteSesion(), []byte(idStr), 0644)

	}
	return errors.New("credenciales inválidas")
}

// CloseSesion cierra la sesión eliminando el archivo local de sesión.
func CloseSesion() {
	os.Remove(getRouteSesion())
}

// GetActiveUser devuelve el ID del usuario con sesión activa.
func GetActiveUser() (string, error) {
	data, err := os.ReadFile(getRouteSesion())
	if err != nil {
		return "", errors.New("no hay una sesión activa. Usa 'task-cli auth login'")
	}
	return strings.TrimSpace(string(data)), nil
}

// ReadPass lee una contraseña sin mostrarla en la terminal.
func ReadPass() string {
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePassword)
}
