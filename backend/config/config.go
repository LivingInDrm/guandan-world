package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port       string
	GinMode    string
	TLSEnabled bool
	TLSCert    string
	TLSKey     string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	AccessSecret       string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func Load() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		JWT:      loadJWTConfig(),
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:       getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
		TLSEnabled: getEnvBool("TLS_ENABLED", false),
		TLSCert:    getEnv("TLS_CERT_PATH", ""),
		TLSKey:     getEnv("TLS_KEY_PATH", ""),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvInt("DB_PORT", 5432),
		Name:            getEnv("DB_NAME", "guandan"),
		User:            getEnv("DB_USER", "guandan"),
		Password:        getEnv("DB_PASSWORD", "guandan"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 5)) * time.Minute,
	}
}

func loadJWTConfig() JWTConfig {
	ginMode := getEnv("GIN_MODE", "debug")
	
	accessSecret := getEnv("JWT_ACCESS_SECRET", "")
	refreshSecret := getEnv("JWT_REFRESH_SECRET", "")
	
	if ginMode == "release" {
		if accessSecret == "" || refreshSecret == "" {
			panic("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set in production mode")
		}
	} else {
		if accessSecret == "" {
			accessSecret = "dev-access-secret-change-in-production"
		}
		if refreshSecret == "" {
			refreshSecret = "dev-refresh-secret-change-in-production"
		}
	}
	
	return JWTConfig{
		AccessSecret:       accessSecret,
		RefreshSecret:      refreshSecret,
		AccessTokenExpiry:  time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_MINUTES", 15)) * time.Minute,
		RefreshTokenExpiry: time.Duration(getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7)) * 24 * time.Hour,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
