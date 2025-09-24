package common

import (
	"fmt"
	"net/http"
)

// ErrorCode représente un code d'erreur
type ErrorCode string

const (
	// Erreurs générales
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeConflict     ErrorCode = "CONFLICT"

	// Erreurs spécifiques aux utilisateurs
	ErrCodeUserNotFound   ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserExists     ErrorCode = "USER_EXISTS"
	ErrCodeInvalidSession ErrorCode = "INVALID_SESSION"
	ErrCodeSessionExpired ErrorCode = "SESSION_EXPIRED"

	// Erreurs spécifiques aux clients
	ErrCodeCustomerNotFound ErrorCode = "CUSTOMER_NOT_FOUND"
	ErrCodeCustomerExists   ErrorCode = "CUSTOMER_EXISTS"

	// Erreurs Ory
	ErrCodeKratosError ErrorCode = "KRATOS_ERROR"
	ErrCodeHydraError  ErrorCode = "HYDRA_ERROR"
	ErrCodeKetoError   ErrorCode = "KETO_ERROR"
)

// AppError représente une erreur de l'application
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
}

// Error implémente l'interface error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError crée une nouvelle erreur d'application
func NewAppError(code ErrorCode, message string, details ...string) *AppError {
	httpStatus := getHTTPStatus(code)
	appErr := &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
	if len(details) > 0 {
		appErr.Details = details[0]
	}
	return appErr
}

// getHTTPStatus retourne le code HTTP correspondant au code d'erreur
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidInput:
		return http.StatusBadRequest
	case ErrCodeNotFound, ErrCodeUserNotFound, ErrCodeCustomerNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized, ErrCodeInvalidSession, ErrCodeSessionExpired:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeConflict, ErrCodeUserExists, ErrCodeCustomerExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Erreurs prédéfinies
var (
	ErrInternalServer = NewAppError(ErrCodeInternal, "Erreur interne du serveur")
	ErrInvalidInput   = NewAppError(ErrCodeInvalidInput, "Données d'entrée invalides")
	ErrNotFound       = NewAppError(ErrCodeNotFound, "Ressource non trouvée")
	ErrUnauthorized   = NewAppError(ErrCodeUnauthorized, "Non autorisé")
	ErrForbidden      = NewAppError(ErrCodeForbidden, "Accès interdit")
	ErrConflict       = NewAppError(ErrCodeConflict, "Conflit de ressources")

	// Erreurs utilisateurs
	ErrUserNotFound   = NewAppError(ErrCodeUserNotFound, "Utilisateur non trouvé")
	ErrUserExists     = NewAppError(ErrCodeUserExists, "Utilisateur déjà existant")
	ErrInvalidSession = NewAppError(ErrCodeInvalidSession, "Session invalide")
	ErrSessionExpired = NewAppError(ErrCodeSessionExpired, "Session expirée")

	// Erreurs clients
	ErrCustomerNotFound = NewAppError(ErrCodeCustomerNotFound, "Client non trouvé")
	ErrCustomerExists   = NewAppError(ErrCodeCustomerExists, "Client déjà existant")

	// Erreurs Ory
	ErrKratosError = NewAppError(ErrCodeKratosError, "Erreur Kratos")
	ErrHydraError  = NewAppError(ErrCodeHydraError, "Erreur Hydra")
	ErrKetoError   = NewAppError(ErrCodeKetoError, "Erreur Keto")
)
