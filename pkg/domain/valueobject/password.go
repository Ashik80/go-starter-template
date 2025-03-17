package valueobject

import "unicode"

type Password struct {
	value string
}

func NewPassword(password string) (Password, error) {
	errs := IsStrongPassword(password)
	if len(errs) > 0 {
		return Password{}, &PasswordValidationError{Errors: errs}
	}
	return Password{value: password}, nil
}

func (p Password) ToString() string {
	return p.value
}

type PasswordValidationError struct {
	Errors []string
}

func (p *PasswordValidationError) Error() string {
	return "Password must have minimum length of 8 characters and must contain at least 1 uppercase letter, 1 lowercase letter, 1 digit and 1 symbol"
}

func IsStrongPassword(password string) []string {
	var errs []string
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			hasSpecial = true
		}
	}

	if !hasMinLen {
		errs = append(errs, "Password must be at least 8 characters long")
	}
	if !hasUpper {
		errs = append(errs, "Password must have at least 1 uppercase character")
	}
	if !hasLower {
		errs = append(errs, "Password must have at least 1 lowercase character")
	}
	if !hasNumber {
		errs = append(errs, "Password must have at least 1 digit")
	}
	if !hasSpecial {
		errs = append(errs, "Password must have at least 1 special character or symbol")
	}
	return errs
}
