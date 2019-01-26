package email

import (
	"errors"
	"regexp"
	"strings"
)

func SanitizeEmail(email string) (string, error) {
	email = strings.Replace(email, " ", "", -1)

	if !validEmail(email) {
		return "", errors.New("Not a valid email address: " + email)
	}

	return email, nil
}

func validEmail(text string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(text)
}
