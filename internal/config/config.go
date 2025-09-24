package config

import (
	"os"
	"strconv"
	"time"
)

// Config contient la configuration de l'application
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Ory      OryConfig      `json:"ory"`
	Logging  LoggingConfig  `json:"logging"`
}

// ServerConfig contient la configuration du serveur
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// DatabaseConfig contient la configuration de la base de données
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

// OryConfig contient la configuration des services Ory
type OryConfig struct {
	Kratos KratosConfig `json:"kratos"`
	Hydra  HydraConfig  `json:"hydra"`
	Keto   KetoConfig   `json:"keto"`
}

// KratosConfig contient la configuration de Kratos
type KratosConfig struct {
	PublicURL string `json:"public_url"`
	AdminURL  string `json:"admin_url"`
}

// HydraConfig contient la configuration de Hydra
type HydraConfig struct {
	PublicURL string `json:"public_url"`
	AdminURL  string `json:"admin_url"`
}

// KetoConfig contient la configuration de Keto
type KetoConfig struct {
	ReadURL  string `json:"read_url"`
	WriteURL string `json:"write_url"`
}

// LoggingConfig contient la configuration du logging
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// Load charge la configuration depuis les variables d'environnement
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getIntEnv("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "ndugu"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Ory: OryConfig{
			Kratos: KratosConfig{
				PublicURL: getEnv("KRATOS_PUBLIC_URL", "http://localhost:4433"),
				AdminURL:  getEnv("KRATOS_ADMIN_URL", "http://localhost:4434"),
			},
			Hydra: HydraConfig{
				PublicURL: getEnv("HYDRA_PUBLIC_URL", "http://localhost:4444"),
				AdminURL:  getEnv("HYDRA_ADMIN_URL", "http://localhost:4445"),
			},
			Keto: KetoConfig{
				ReadURL:  getEnv("KETO_READ_URL", "http://localhost:4466"),
				WriteURL: getEnv("KETO_WRITE_URL", "http://localhost:4467"),
			},
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}
}

// getEnv récupère une variable d'environnement avec une valeur par défaut
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv récupère une variable d'environnement entière avec une valeur par défaut
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv récupère une variable d'environnement de durée avec une valeur par défaut
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
