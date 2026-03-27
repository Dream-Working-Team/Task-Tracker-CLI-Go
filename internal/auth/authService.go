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

// archiveSesion is the local file that stores the currently authenticated user id.
const archiveSesion = ".session"

// passwordHashCost controls bcrypt hashing strength.
// Higher values improve security but increase processing time.
const passwordHashCost = bcrypt.DefaultCost

// AuthService coordinates user registration, login, and session persistence.
type AuthService struct {
	save *storage.Storage
}

// NewAuthService creates an AuthService with the provided storage backend.
func NewAuthService(s *storage.Storage) *AuthService {
	return &AuthService{save: s}
}

// Register creates a new user with a hashed password if the username is not taken.
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

// Login validates credentials and stores the authenticated user id in the session file.
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

// CloseSesion removes the session file to end the current user session.
func CloseSesion() {
	os.Remove(archiveSesion)
}

// GetActiveUser returns the currently authenticated user id from the session file.
func GetActiveUser() (string, error) {
	data, err := os.ReadFile(archiveSesion)
	if err != nil {
		return "", errors.New("no hay una sesión activa. Usa 'task-cli auth login'")
	}
	return strings.TrimSpace(string(data)), nil
}

// ReadPass reads a password from terminal input without echoing typed characters.
func ReadPass() string {
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePassword)
}
