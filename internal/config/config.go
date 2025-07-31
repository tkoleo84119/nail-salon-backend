package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type DBConfig struct {
	DSN                   string
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	SSLMode               string
	MaxOpenConns          int
	MaxConnMaxLife        time.Duration
	MaxConnLifetimeJitter time.Duration
	MaxConnIdleTime       time.Duration
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type LineConfig struct {
	ChannelID        string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type ServerConfig struct {
	Port            string
	SnowflakeNodeId int64
}

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Line   LineConfig
	Server ServerConfig
	CORS   CORSConfig
}

func Load() *Config {
	dbConfig := DBConfig{
		DSN:                   os.Getenv("DB_DSN"),
		Host:                  getenvDefault("DB_HOST", "localhost"),
		Port:                  getenvIntDefault("DB_PORT", 5432),
		User:                  os.Getenv("DB_USER"),
		Password:              os.Getenv("DB_PASS"),
		Name:                  os.Getenv("DB_NAME"),
		SSLMode:               getenvDefault("DB_SSLMODE", "disable"),
		MaxOpenConns:          getenvIntDefault("DB_MAX_OPEN_CONNS", 25),
		MaxConnMaxLife:        getenvDuration("DB_CONN_MAX_LIFE", "30m"),
		MaxConnLifetimeJitter: getenvDuration("DB_CONN_MAX_LIFE_JITTER", "5m"),
		MaxConnIdleTime:       getenvDuration("DB_CONN_IDLE_TIME", "2m"),
	}

	if dbConfig.DSN == "" {
		dbConfig.DSN = dbConfig.ConnectionString()
	}

	jwtConfig := JWTConfig{
		Secret:      getenvRequired("JWT_SECRET"),
		ExpiryHours: getenvIntDefault("JWT_EXPIRY_HOURS", 1),
	}

	lineConfig := LineConfig{
		ChannelID: getenvRequired("LINE_CHANNEL_ID"),
	}

	serverConfig := ServerConfig{
		Port:            getenvDefault("PORT", "3000"),
		SnowflakeNodeId: int64(getenvIntDefault("SNOWFLAKE_NODE_ID", 1)),
	}

	corsConfig := CORSConfig{
		AllowedOrigins:   getenvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		AllowedMethods:   getenvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
		AllowedHeaders:   getenvSlice("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"}),
		ExposedHeaders:   getenvSlice("CORS_EXPOSED_HEADERS", []string{}),
		AllowCredentials: getenvBoolDefault("CORS_ALLOW_CREDENTIALS", true),
		MaxAge:           getenvIntDefault("CORS_MAX_AGE", 300),
	}

	return &Config{
		DB:     dbConfig,
		JWT:    jwtConfig,
		Line:   lineConfig,
		Server: serverConfig,
		CORS:   corsConfig,
	}
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

func getenvRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return val
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

func getenvSlice(key string, defaultVal []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	
	// Split by comma and trim spaces
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	if len(result) == 0 {
		return defaultVal
	}
	
	return result
}

func getenvBoolDefault(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	
	return parsed
}
