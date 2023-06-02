package utils

import (
	"net/mail"
)

func Is_email_valid(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}
