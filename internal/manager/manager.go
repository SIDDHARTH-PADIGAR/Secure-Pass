package manager

import (
	"fmt"
	"time"

	"secure-pass/internal/generator"
	"secure-pass/internal/storage"
	"secure-pass/internal/utils"

	"github.com/atotto/clipboard"
)

// PasswordManager handles the main logic of the password manager
type PasswordManager struct {
	masterKey []byte
}

// NewPasswordManager creates a new instance of PasswordManager
func NewPasswordManager(key []byte) *PasswordManager {
	return &PasswordManager{
		masterKey: key,
	}
}

// RunMainLoop runs the main menu loop for the password manager
func (pm *PasswordManager) RunMainLoop() {
	for {
		fmt.Println("\nPassword Manager Menu:")
		fmt.Println("1. Save Password")
		fmt.Println("2. Get Password History")
		fmt.Println("3. Generate Strong Password")
		fmt.Println("4. Search Passwords")
		fmt.Println("5. Check Password Expiry Status")
		fmt.Println("6. Copy Password to Clipboard")
		fmt.Println("7. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			pm.handleSavePassword()
		case 2:
			pm.handleGetPasswordHistory()
		case 3:
			pm.handleGeneratePassword()
		case 4:
			pm.handleSearchPasswords()
		case 5:
			pm.handleCheckPasswordExpiry()
		case 6:
			pm.handleCopyToClipboard()
		case 7:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

// handleSavePassword handles the save password option
func (pm *PasswordManager) handleSavePassword() {
	website := utils.ReadInput("Enter Website/App Name: ")
	account := utils.ReadInput("Enter Account Name: ")
	password := utils.ReadSecureInput("Enter Password: ")
	expiryDays := storage.DefaultExpiry

	if utils.GetYesNo("Set custom expiry period? (y/n): ") {
		fmt.Print("Enter number of days until expiry: ")
		fmt.Scanln(&expiryDays)
	}

	if err := storage.SavePassword(website, account, password, pm.masterKey, expiryDays); err != nil {
		fmt.Printf("Error saving password: %v\n", err)
		return
	}
	fmt.Println("Password saved successfully!")
}

// handleGetPasswordHistory handles the get password history option
func (pm *PasswordManager) handleGetPasswordHistory() {
	website := utils.ReadInput("Enter Website/App Name: ")
	account := utils.ReadInput("Enter Account Name: ")

	entries, err := storage.GetPasswordHistory(website, account, pm.masterKey)
	if err != nil {
		fmt.Printf("Error retrieving password history: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No password history found.")
		return
	}

	fmt.Printf("\nPassword History for %s - %s\n", website, account)
	fmt.Println("---------------------------------")
	for _, entry := range entries {
		fmt.Printf("Password: %s | Saved on: %s | Expires in: %d days\n",
			entry.Password, entry.Timestamp, entry.ExpiryDays)
	}
}

// handleGeneratePassword handles the generate password option
func (pm *PasswordManager) handleGeneratePassword() {
	opts := generator.GetPasswordGeneratorOptions()
	password, err := generator.GeneratePassword(opts)
	if err != nil {
		fmt.Printf("Error generating password: %v\n", err)
		return
	}

	fmt.Printf("\nGenerated Password: %s\n", password)

	if utils.GetYesNo("\nWould you like to save this password? (y/n): ") {
		website := utils.ReadInput("Enter Website/App Name: ")
		account := utils.ReadInput("Enter Account Name: ")
		expiryDays := storage.DefaultExpiry

		if utils.GetYesNo("Set custom expiry period? (y/n): ") {
			fmt.Print("Enter number of days until expiry: ")
			fmt.Scanln(&expiryDays)
		}

		if err := storage.SavePassword(website, account, password, pm.masterKey, expiryDays); err != nil {
			fmt.Printf("Error saving password: %v\n", err)
			return
		}
		fmt.Println("Password saved successfully!")
	}
}

// handleSearchPasswords handles the search passwords option
func (pm *PasswordManager) handleSearchPasswords() {
	searchTerm := utils.ReadInput("Enter search term: ")
	results, err := storage.SearchPasswords(searchTerm, pm.masterKey)
	if err != nil {
		fmt.Printf("Error searching passwords: %v\n", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	fmt.Println("\nSearch Results:")
	fmt.Println("--------------")
	for i, result := range results {
		fmt.Printf("%d. Website: %s, Account: %s\n", i+1, result.Website, result.Account)
		isExpiring, daysLeft := storage.CheckPasswordExpiry(result.Entry.Timestamp, result.Entry.ExpiryDays)
		if isExpiring {
			fmt.Printf("   ⚠️ Password expires in %d days!\n", daysLeft)
		}
	}
}

// handleCheckPasswordExpiry handles the check password expiry option
func (pm *PasswordManager) handleCheckPasswordExpiry() {
	website := utils.ReadInput("Enter Website/App Name: ")
	account := utils.ReadInput("Enter Account Name: ")

	entries, err := storage.GetPasswordHistory(website, account, pm.masterKey)
	if err != nil || len(entries) == 0 {
		fmt.Println("No passwords found.")
		return
	}

	latestEntry := entries[len(entries)-1]
	isExpiring, daysLeft := storage.CheckPasswordExpiry(latestEntry.Timestamp, latestEntry.ExpiryDays)

	if isExpiring {
		fmt.Printf("⚠️ WARNING: Password will expire in %d days!\n", daysLeft)
	} else {
		fmt.Printf("✅ Password is valid for %d more days.\n", daysLeft)
	}
}

// handleCopyToClipboard handles the copy password to clipboard option
func (pm *PasswordManager) handleCopyToClipboard() {
	website := utils.ReadInput("Enter Website/App Name: ")
	account := utils.ReadInput("Enter Account Name: ")

	entries, err := storage.GetPasswordHistory(website, account, pm.masterKey)
	if err != nil || len(entries) == 0 {
		fmt.Println("No passwords found.")
		return
	}

	password := entries[len(entries)-1].Password
	if err := copyToClipboard(password); err != nil {
		fmt.Printf("Error copying to clipboard: %v\n", err)
		return
	}

	fmt.Println("🔑 Password copied to clipboard! (Will clear in 30 seconds)")
}

// copyToClipboard copies the password to clipboard with auto-clear after 30 seconds
func copyToClipboard(password string) error {
	err := clipboard.WriteAll(password)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}
	// Clear clipboard after 30 seconds for security
	go func() {
		time.Sleep(30 * time.Second)
		clipboard.WriteAll("")
	}()
	return nil
}
