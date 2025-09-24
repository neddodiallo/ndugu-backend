package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

// SugarLogger implémente l'interface Logger avec logrus
type SugarLogger struct {
	logger *logrus.Logger
}

// NewSugarLogger crée une nouvelle instance de SugarLogger
func NewSugarLogger() Logger {
	logger := logrus.New()

	// Configuration du format
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	// Configuration du niveau de log
	logger.SetLevel(logrus.InfoLevel)

	// Configuration de la sortie
	logger.SetOutput(os.Stdout)

	return &SugarLogger{
		logger: logger,
	}
}

// Info enregistre un message d'information
func (l *SugarLogger) Info(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn enregistre un message d'avertissement
func (l *SugarLogger) Warn(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error enregistre un message d'erreur
func (l *SugarLogger) Error(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Debug enregistre un message de débogage
func (l *SugarLogger) Debug(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// WithField ajoute un champ au logger
func (l *SugarLogger) WithField(key string, value interface{}) Logger {
	return &SugarLogger{
		logger: l.logger.WithField(key, value).Logger,
	}
}

// WithFields ajoute plusieurs champs au logger
func (l *SugarLogger) WithFields(fields map[string]interface{}) Logger {
	return &SugarLogger{
		logger: l.logger.WithFields(fields).Logger,
	}
}

// SetLevel définit le niveau de log
func (l *SugarLogger) SetLevel(level string) {
	switch level {
	case "debug":
		l.logger.SetLevel(logrus.DebugLevel)
	case "info":
		l.logger.SetLevel(logrus.InfoLevel)
	case "warn":
		l.logger.SetLevel(logrus.WarnLevel)
	case "error":
		l.logger.SetLevel(logrus.ErrorLevel)
	default:
		l.logger.SetLevel(logrus.InfoLevel)
	}
}
