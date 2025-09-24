package common

import (
	"log"
	"os"
)

// Logger interface pour le logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	SetLevel(level string)
}

// SimpleLogger implémentation simple du logger
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

// NewSimpleLogger crée un nouveau logger simple
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info log un message d'information
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLogger.Printf(msg, args...)
	} else {
		l.infoLogger.Println(msg)
	}
}

// Error log un message d'erreur
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.errorLogger.Printf(msg, args...)
	} else {
		l.errorLogger.Println(msg)
	}
}

// Debug log un message de debug
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.debugLogger.Printf(msg, args...)
	} else {
		l.debugLogger.Println(msg)
	}
}

// Warn log un message d'avertissement
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.warnLogger.Printf(msg, args...)
	} else {
		l.warnLogger.Println(msg)
	}
}

// WithField ajoute un champ au logger (implémentation simple)
func (l *SimpleLogger) WithField(key string, value interface{}) Logger {
	// Pour SimpleLogger, on retourne le même logger
	// Les champs ne sont pas supportés dans cette implémentation simple
	return l
}

// WithFields ajoute plusieurs champs au logger (implémentation simple)
func (l *SimpleLogger) WithFields(fields map[string]interface{}) Logger {
	// Pour SimpleLogger, on retourne le même logger
	// Les champs ne sont pas supportés dans cette implémentation simple
	return l
}

// SetLevel définit le niveau de log (implémentation simple)
func (l *SimpleLogger) SetLevel(level string) {
	// SimpleLogger ne supporte pas la configuration de niveau
	// Cette méthode est présente pour l'interface
}

// Logger global
var DefaultLogger Logger = NewSimpleLogger()

// Fonctions de convenance
func Info(msg string, args ...interface{}) {
	DefaultLogger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	DefaultLogger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	DefaultLogger.Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	DefaultLogger.Warn(msg, args...)
}
