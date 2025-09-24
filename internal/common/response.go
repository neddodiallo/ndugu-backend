package common

import (
	"encoding/json"
	"net/http"
)

// Response représente une réponse standard de l'API
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *AppError   `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse crée une réponse de succès
func SuccessResponse(data interface{}, message ...string) *Response {
	resp := &Response{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return resp
}

// ErrorResponse crée une réponse d'erreur
func ErrorResponse(err *AppError) *Response {
	return &Response{
		Success: false,
		Error:   err,
		Message: err.Message,
	}
}

// WriteJSON écrit une réponse JSON
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccess écrit une réponse de succès
func WriteSuccess(w http.ResponseWriter, data interface{}, message ...string) {
	resp := SuccessResponse(data, message...)
	WriteJSON(w, http.StatusOK, resp)
}

// WriteError écrit une réponse d'erreur
func WriteError(w http.ResponseWriter, err *AppError) {
	resp := ErrorResponse(err)
	WriteJSON(w, err.HTTPStatus, resp)
}

// WriteInternalError écrit une erreur interne
func WriteInternalError(w http.ResponseWriter, message string) {
	err := NewAppError(ErrCodeInternal, message)
	WriteError(w, err)
}
