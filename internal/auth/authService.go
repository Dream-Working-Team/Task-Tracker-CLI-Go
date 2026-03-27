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
)

// archiveSesion es el archivo local donde se corre la sesion actual del usuario
const archiveSesion = ".session"

// passwordHashCost controla la fuerza del hash de bcrypt.
// Valores más altos mejoran la seguridad, pero aumentan el tiempo de procesamiento.
const passwordHashCost = bcrypt.DefaultCost

// AuthService coordina el registro de usuarios, el inicio de sesión y la persistencia de sesión.
type AuthService struct {
	save *storage.Storage
}

// NewAuthService crea un AuthService con el backend de almacenamiento proporcionado.
func NewAuthService(s *storage.Storage) *AuthService {
	return &AuthService{save: s}
}

// Register crea un nuevo usuario con contraseña hasheada si el nombre de usuario no está en uso.
func (s *AuthService) Register(username, password string) error {
	usersGeneral, err := s.save.ReadUsers()
	if err != nil {
		return err
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

	usersGeneral = append(usersGeneral, model.User{Username: username, PasswordHash: string(hash)})
	return s.save.SaveUsers(usersGeneral)
}

// Login valida las credenciales y guarda el id del usuario autenticado en el archivo de sesión.
func (s *AuthService) Login(username, password string) error {
	usersGeneral, err := s.save.ReadUsers()
	if err != nil {
		return err
	}

	for _, u := range usersGeneral {
		if u.Username == username {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
			return errors.New("credenciales inválidas")
		}
		idStr := strconv.Itoa(u.ID)
		return os.WriteFile(archiveSesion, []byte(idStr), 0644)

	}
	return errors.New("credenciales inválidas")
}

// CloseSesion elimina el archivo de sesión para cerrar la sesión del usuario actual.
func CloseSesion() {
	os.Remove(archiveSesion)
}

// GetActiveUser devuelve el id del usuario autenticado actualmente desde el archivo de sesión.
func GetActiveUser() (string, error) {
	data, err := os.ReadFile(archiveSesion)
	if err != nil {
		return "", errors.New("no hay una sesión activa. Usa 'task-cli auth login'")
	}
	return strings.TrimSpace(string(data)), nil
}

// ReadPass lee una contraseña desde la terminal sin mostrar los caracteres escritos.
func ReadPass() string {
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePassword)
}
