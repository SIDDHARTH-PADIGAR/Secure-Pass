package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadSecureInput reads input from the console without echoing it
func ReadSecureInput(prompt string) string {
	fmt.Print(prompt)
	password, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(password)
}

// ReadInput reads a line of input from the console
func ReadInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// GetYesNo prompts for a yes/no answer and returns true for yes
func GetYesNo(prompt string) bool {
	var response string
	fmt.Print(prompt)
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

// RemoveChars removes specified characters from a string
func RemoveChars(str, chars string) string {
	for _, c := range chars {
		str = strings.ReplaceAll(str, string(c), "")
	}
	return str
}
