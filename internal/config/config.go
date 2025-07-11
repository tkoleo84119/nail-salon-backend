package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type DBConfig struct {
	DSN          string
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
}

func LoadString() DBConfig {
	dbConfig := DBConfig{
		DSN:          os.Getenv("DB_DSN"),
		Host:         getenvDefault("DB_HOST", "localhost"),
		Port:         getenvIntDefault("DB_PORT", 5432),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASS"),
		Name:         os.Getenv("DB_NAME"),
		SSLMode:      getenvDefault("DB_SSLMODE", "disable"),
		MaxOpenConns: getenvIntDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getenvIntDefault("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLife:  getenvDuration("DB_CONN_MAX_LIFE", "30m"),
	}

	if dbConfig.DSN == "" {
		dbConfig.DSN = dbConfig.ConnectionString()
	}

	return dbConfig
}

func (c DBConfig) ConnectionString() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func getenvDefault(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func getenvIntDefault(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return defaultVal
		}
		return parsed
	}
	return defaultVal
}

func getenvDuration(key, defaultStr string) time.Duration {
	raw := os.Getenv(key)

	// default value must be able to parse successfully; otherwise, panic
	def, err := time.ParseDuration(defaultStr)
	if err != nil {
		panic("config: invalid default duration for " + key + ": " + err.Error())
	}

	if raw == "" {
		return def
	}

	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("[WARN] config: %s=%q is not a valid duration, fallback to %q: %v", key, raw, defaultStr, err)
		return def
	}

	return d
}
