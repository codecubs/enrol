package email

import "testing"

func TestSanitizeEmail(t *testing.T) {
	input := "Hello"
	_, err := SanitizeEmail(input)
	if err == nil {
		t.Error(`"Hello" is not a valid email`)
	}
}
