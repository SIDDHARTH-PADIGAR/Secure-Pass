package main

import (
	"log"

	"secure-pass/internal/auth"
	"secure-pass/internal/manager"
	"secure-pass/internal/utils"
)

func main() {
	if err := auth.InitializeMasterPassword(); err != nil {
		log.Fatalf("Failed to initialize master password: %v", err)
	}

	masterPassword := utils.ReadSecureInput("Enter Master Password: ")
	key, err := auth.VerifyMasterPassword(masterPassword)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	passwordManager := manager.NewPasswordManager(key)
	passwordManager.RunMainLoop()
}
