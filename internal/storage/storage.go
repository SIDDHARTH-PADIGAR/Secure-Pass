package storage

import (
	"fmt"
	"strings"
	"time"

	"secure-pass/internal/crypto"

	"github.com/tidwall/buntdb"
)

const (
	DefaultExpiry = 90 // Default password expiry in days
)

type PasswordEntry struct {
	Password   string
	Timestamp  string
	ExpiryDays int    // Number of days until password expires
	LastCopied string // Timestamp of last clipboard copy
}

// SavePassword saves or updates a password in the database
func SavePassword(website, account, password string, key []byte, expiryDays int) error {
	if expiryDays <= 0 {
		expiryDays = DefaultExpiry
	}

	db, err := buntdb.Open("securepass.db")
	if err != nil {
		return err
	}
	defer db.Close()

	encryptedPassword, err := crypto.Encrypt(password, key)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("%s | %s | %d", encryptedPassword, timestamp, expiryDays)

	return db.Update(func(tx *buntdb.Tx) error {
		key := fmt.Sprintf("%s:%s", website, account)
		val, err := tx.Get(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}

		var newHistory string
		if err == buntdb.ErrNotFound {
			newHistory = entry
		} else {
			newHistory = strings.TrimSpace(val) + "\n" + entry
		}

		_, _, err = tx.Set(key, newHistory, nil)
		return err
	})
}

// GetPasswordHistory retrieves password history for a website and account
func GetPasswordHistory(website, account string, key []byte) ([]PasswordEntry, error) {
	db, err := buntdb.Open("securepass.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var entries []PasswordEntry
	err = db.View(func(tx *buntdb.Tx) error {
		history, err := tx.Get(fmt.Sprintf("%s:%s", website, account))
		if err == buntdb.ErrNotFound {
			return nil
		}
		if err != nil {
			return err
		}

		for _, line := range strings.Split(strings.TrimSpace(history), "\n") {
			parts := strings.Split(line, " | ")
			if len(parts) < 2 {
				continue
			}

			decryptedPassword, err := crypto.Decrypt(strings.TrimSpace(parts[0]), key)
			if err != nil {
				continue
			}

			expiryDays := DefaultExpiry
			if len(parts) >= 3 {
				fmt.Sscanf(parts[2], "%d", &expiryDays)
			}

			entries = append(entries, PasswordEntry{
				Password:   decryptedPassword,
				Timestamp:  strings.TrimSpace(parts[1]),
				ExpiryDays: expiryDays,
			})
		}
		return nil
	})

	return entries, err
}

// SearchPasswords searches for passwords matching the searchTerm
func SearchPasswords(searchTerm string, masterKey []byte) ([]struct {
	Website string
	Account string
	Entry   PasswordEntry
}, error) {
	db, err := buntdb.Open("securepass.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var results []struct {
		Website string
		Account string
		Entry   PasswordEntry
	}

	err = db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(dbKey, value string) bool {
			website, account, found := strings.Cut(dbKey, ":")
			if !found {
				return true
			}

			// Case-insensitive search
			if strings.Contains(strings.ToLower(website), strings.ToLower(searchTerm)) ||
				strings.Contains(strings.ToLower(account), strings.ToLower(searchTerm)) {

				entries, err := GetPasswordHistory(website, account, masterKey)
				if err != nil || len(entries) == 0 {
					return true
				}

				results = append(results, struct {
					Website string
					Account string
					Entry   PasswordEntry
				}{
					Website: website,
					Account: account,
					Entry:   entries[len(entries)-1], // Most recent entry
				})
			}
			return true
		})
	})

	return results, err
}

// CheckPasswordExpiry checks if a password is close to expiry
func CheckPasswordExpiry(timestamp string, expiryDays int) (bool, int) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		return false, 0
	}

	expiryDate := parsedTime.AddDate(0, 0, expiryDays)
	daysUntilExpiry := int(time.Until(expiryDate).Hours() / 24)

	return daysUntilExpiry <= 7, daysUntilExpiry // Alert if 7 or fewer days until expiry
}
