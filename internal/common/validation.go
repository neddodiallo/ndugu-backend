package common

import (
	"regexp"
	"strings"
)

// Validator interface pour la validation
type Validator interface {
	Validate() error
}

// EmailRegex pour valider les emails
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// PhoneRegex pour valider les numéros de téléphone
var PhoneRegex = regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)

// ValidateEmail valide un email
func ValidateEmail(email string) error {
	if email == "" {
		return NewAppError(ErrCodeInvalidInput, "Email requis")
	}
	if !EmailRegex.MatchString(email) {
		return NewAppError(ErrCodeInvalidInput, "Format d'email invalide")
	}
	return nil
}

// ValidatePhone valide un numéro de téléphone
func ValidatePhone(phone string) error {
	if phone == "" {
		return NewAppError(ErrCodeInvalidInput, "Numéro de téléphone requis")
	}
	if !PhoneRegex.MatchString(phone) {
		return NewAppError(ErrCodeInvalidInput, "Format de numéro de téléphone invalide")
	}
	return nil
}

// ValidatePassword valide un mot de passe
func ValidatePassword(password string) error {
	if password == "" {
		return NewAppError(ErrCodeInvalidInput, "Mot de passe requis")
	}
	if len(password) < 8 {
		return NewAppError(ErrCodeInvalidInput, "Le mot de passe doit contenir au moins 8 caractères")
	}
	return nil
}

// ValidateRequired valide qu'un champ est requis
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return NewAppError(ErrCodeInvalidInput, fieldName+" est requis")
	}
	return nil
}

// ValidateLength valide la longueur d'une chaîne
func ValidateLength(value, fieldName string, min, max int) error {
	length := len(strings.TrimSpace(value))
	if length < min {
		return NewAppError(ErrCodeInvalidInput, fieldName+" doit contenir au moins "+string(rune(min))+" caractères")
	}
	if length > max {
		return NewAppError(ErrCodeInvalidInput, fieldName+" ne peut pas dépasser "+string(rune(max))+" caractères")
	}
	return nil
}
