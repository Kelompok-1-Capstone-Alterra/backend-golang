package utils

import "regexp"

func IsValidPhone(phone string) bool {
	// Regular expression to match phone number with country code
	// Modify this pattern according to your desired format
	pattern := `^\+[1-9]\d{1,14}$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Check if the phone matches the pattern
	return regex.MatchString(phone)
}
