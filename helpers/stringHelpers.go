package helpers

import (
	"net/mail"
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

func IsPasswordStrong(pass string) FieldError {
	var letters = 0
	var number, upper, special bool
	for _, c := range pass {
		switch {
		case unicode.IsNumber(c):
			number = true
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
	if letters < 7 {
		return FieldError{
			Message: "Password must be at least 7 characters",
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
