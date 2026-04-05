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

// AuthService coordinates user registration, login, and session persistence
type AuthService struct {
	save *storage.Storage
}

// getSessionPath returns the local file path where the active session is stored
func getSessionPath() string {
	dir, err := config.GetDataDirectory()
	if err != nil {
		return ".ActiveSession" // Emergency fallback to the local root directory
	}
	return filepath.Join(dir, ".ActiveSession")
}

// NewAuthService creates a new authentication service instance
func NewAuthService(s *storage.Storage) *AuthService {
	return &AuthService{save: s}
}

// Register creates a new user if the username does not already exist
func (s *AuthService) Register(username, password string) error {
	users, err := s.save.ReadUsers()
	if err != nil {
		return err
	}

	newID := 1
	if len(users) > 0 {
		newID = users[len(users)-1].ID + 1
	}

	for _, u := range users {
		if u.Username == username {
			return errors.New("The username already exists")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users = append(users, model.User{
		ID:           newID,
		Username:     username,
		PasswordHash: string(hash),
	})
	return s.save.SaveUsers(users)
}

// Login validates credentials and stores the authenticated user session
func (s *AuthService) Login(username, password string) error {
	users, err := s.save.ReadUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Username != username {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
			return errors.New("Invalid credentials")
		}
		idStr := strconv.Itoa(u.ID)
		return os.WriteFile(getSessionPath(), []byte(idStr), 0644)

	}
	return errors.New("Invalid credentials")
}

// CloseSession closes the session by removing the local session file
func CloseSession() {
	os.Remove(getSessionPath())
}

// GetActiveUser returns the ID of the user with an active session
func GetActiveUser() (string, error) {
	data, err := os.ReadFile(getSessionPath())
	if err != nil {
		return "", errors.New("No active session. Use 'task-cli auth login'")
	}
	return strings.TrimSpace(string(data)), nil
}

// ReadPass reads a password without displaying it in the terminal
func ReadPass() string {
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePassword)
}
