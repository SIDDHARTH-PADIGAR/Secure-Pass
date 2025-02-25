package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"secure-pass/internal/utils"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 32
	iterations = 100000
	keySize    = 32
	masterFile = "master.key"
)

// InitializeMasterPassword initializes or verifies the master password
func InitializeMasterPassword() error {
	if _, err := os.Stat(masterFile); os.IsNotExist(err) {
		return createMasterPassword()
	}
	return nil
}

// createMasterPassword creates a new master password
func createMasterPassword() error {
	fmt.Println("First time setup - Create Master Password")
	password := utils.ReadSecureInput("Enter Master Password: ")
	confirmPassword := utils.ReadSecureInput("Confirm Master Password: ")

	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)

	file, err := os.OpenFile(masterFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to create master key file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(salt); err != nil {
		return fmt.Errorf("failed to write salt: %w", err)
	}
	if _, err := file.Write(hash); err != nil {
		return fmt.Errorf("failed to write hash: %w", err)
	}

	return nil
}

// VerifyMasterPassword verifies the master password and returns the encryption key
func VerifyMasterPassword(password string) ([]byte, error) {
	file, err := os.ReadFile(masterFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read master key file: %w", err)
	}

	salt := file[:saltSize]
	storedHash := file[saltSize:]
	hash := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)

	if !compareHashes(hash, storedHash) {
		return nil, errors.New("incorrect master password")
	}

	return hash, nil
}

// compareHashes compares hashes in constant time
func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
