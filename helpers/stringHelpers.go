package helpers

import (
	"net/mail"
	"strings"
	"unicode"
)

type FieldError struct {
	Message string
	Success bool
}

func IsEmpty(data string) bool {
	if len(data) == 0 {
		return true
	}
	return false
}

func IsEmail(data string) bool {
	_, err := mail.ParseAddress(data)
	return err == nil

}

// GetStringInBetween Returns empty string if no start string found
func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

func IsPasswordStrong(pass string) FieldError {
	var letters = 0
	var number, upper, special bool
	for _, c := range pass {
		switch {
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsSymbol(c) || unicode.IsPunct(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
		}
	}
	if letters < 6 {
		return FieldError{
			Message: "Password must be at least 6 characters",
			Success: false,
		}
	} else if !number {
		return FieldError{
			Message: "Password must have at least 1 number",
			Success: false,
		}
	} else if !upper {
		return FieldError{
			Message: "Password must have at least 1 uppercase",
			Success: false,
		}
	} else if !special {
		return FieldError{
			Message: "Password must have special character",
			Success: false,
		}
	} else {
		return FieldError{
			Success: true,
		}
	}
}
