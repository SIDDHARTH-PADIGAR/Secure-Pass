package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"secure-pass/internal/utils"
)

const (
	upperChars     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars     = "abcdefghijklmnopqrstuvwxyz"
	numberChars    = "0123456789"
	symbolChars    = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	ambiguousChars = "Il1O0"
)

type PasswordOptions struct {
	Length           int
	IncludeUpper     bool
	IncludeLower     bool
	IncludeNumbers   bool
	IncludeSymbols   bool
	ExcludeAmbiguous bool
}

// GeneratePassword generates a strong password based on the given options
func GeneratePassword(opts PasswordOptions) (string, error) {
	if opts.Length < 8 {
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	var charset strings.Builder
	var mandatoryChars []rune

	if opts.IncludeLower {
		charset.WriteString(lowerChars)
		mandatoryChars = append(mandatoryChars, rune(lowerChars[0]))
	}
	if opts.IncludeUpper {
		charset.WriteString(upperChars)
		mandatoryChars = append(mandatoryChars, rune(upperChars[0]))
	}
	if opts.IncludeNumbers {
		charset.WriteString(numberChars)
		mandatoryChars = append(mandatoryChars, rune(numberChars[0]))
	}
	if opts.IncludeSymbols {
		charset.WriteString(symbolChars)
		mandatoryChars = append(mandatoryChars, rune(symbolChars[0]))
	}

	if charset.Len() == 0 {
		return "", fmt.Errorf("at least one character type must be selected")
	}

	chars := charset.String()
	if opts.ExcludeAmbiguous {
		chars = utils.RemoveChars(chars, ambiguousChars)
	}

	var password strings.Builder
	password.Grow(opts.Length)

	for i := 0; i < len(mandatoryChars) && i < opts.Length; i++ {
		password.WriteRune(mandatoryChars[i])
	}

	for i := password.Len(); i < opts.Length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password.WriteByte(chars[randomIndex.Int64()])
	}

	passwordRunes := []rune(password.String())
	for i := len(passwordRunes) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", fmt.Errorf("failed to shuffle password: %w", err)
		}
		passwordRunes[i], passwordRunes[j.Int64()] = passwordRunes[j.Int64()], passwordRunes[i]
	}

	return string(passwordRunes), nil
}

// GetPasswordGeneratorOptions prompts the user for password generation options
func GetPasswordGeneratorOptions() PasswordOptions {
	var opts PasswordOptions

	fmt.Println("\nPassword Generator Options:")
	fmt.Print("Enter password length (minimum 8): ")
	fmt.Scanln(&opts.Length)
	if opts.Length < 8 {
		opts.Length = 8
	}

	opts.IncludeLower = utils.GetYesNo("Include lowercase letters? (y/n): ")
	opts.IncludeUpper = utils.GetYesNo("Include uppercase letters? (y/n): ")
	opts.IncludeNumbers = utils.GetYesNo("Include numbers? (y/n): ")
	opts.IncludeSymbols = utils.GetYesNo("Include symbols? (y/n): ")
	opts.ExcludeAmbiguous = utils.GetYesNo("Exclude ambiguous characters (I, l, 1, O, 0)? (y/n): ")

	return opts
}
