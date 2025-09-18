package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
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
	ExpiryHours float64
}

type LineConfig struct {
	LiffChannelID      string
	MessageAccessToken string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type SchedulerConfig struct {
	RefreshRevokeCron string
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

type ProxyConfig struct {
	TrustedProxies []string
}

type CookieConfig struct {
	// Cookie names
	AdminRefreshName    string
	CustomerRefreshName string

	// Common attributes
	Domain   string
	Path     string
	Secure   bool
	SameSite string // one of: Lax, Strict, None, Default

	// Expiration settings (in days)
	AdminRefreshMaxAgeDays    int
	CustomerRefreshMaxAgeDays int
}

type Config struct {
	DB        DBConfig
	JWT       JWTConfig
	Line      LineConfig
	Redis     RedisConfig
	Scheduler SchedulerConfig
	Server    ServerConfig
	CORS      CORSConfig
	Cookie    CookieConfig
	Proxy     ProxyConfig
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
		ExpiryHours: getenvFloatDefault("JWT_EXPIRY_HOURS", 1.0),
	}

	lineConfig := LineConfig{
		LiffChannelID:      getenvRequired("LINE_LIFF_CHANNEL_ID"),
		MessageAccessToken: getenvRequired("LINE_MESSAGING_ACCESS_TOKEN"),
	}

	redisConfig := RedisConfig{
		Host:     getenvDefault("REDIS_HOST", "localhost"),
		Port:     getenvIntDefault("REDIS_PORT", 6379),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       getenvIntDefault("REDIS_DB", 0),
	}

	schedulerConfig := SchedulerConfig{
		RefreshRevokeCron: getAndCheckCronExpression("REFRESH_REVOKE_CRON"),
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

	cookieConfig := CookieConfig{
		AdminRefreshName:          getenvDefault("ADMIN_REFRESH_COOKIE_NAME", "admin_refresh_token"),
		CustomerRefreshName:       getenvDefault("CUSTOMER_REFRESH_COOKIE_NAME", "customer_refresh_token"),
		Domain:                    os.Getenv("COOKIE_DOMAIN"),
		Path:                      getenvDefault("COOKIE_PATH", "/"),
		Secure:                    getenvBoolDefault("COOKIE_SECURE", false),
		SameSite:                  getenvDefault("COOKIE_SAMESITE", "Lax"),
		AdminRefreshMaxAgeDays:    getenvIntDefault("ADMIN_REFRESH_COOKIE_MAX_AGE_DAYS", 7),
		CustomerRefreshMaxAgeDays: getenvIntDefault("CUSTOMER_REFRESH_COOKIE_MAX_AGE_DAYS", 7),
	}

	proxyConfig := ProxyConfig{
		TrustedProxies: getenvSlice("TRUSTED_PROXIES", []string{}),
	}

	return &Config{
		DB:        dbConfig,
		JWT:       jwtConfig,
		Line:      lineConfig,
		Redis:     redisConfig,
		Scheduler: schedulerConfig,
		Server:    serverConfig,
		CORS:      corsConfig,
		Cookie:    cookieConfig,
		Proxy:     proxyConfig,
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

func getenvFloatDefault(key string, defaultVal float64) float64 {
	if val, exists := os.LookupEnv(key); exists {
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return defaultVal
		}
		return parsed
	}
	return defaultVal
}

func getAndCheckCronExpression(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}

	// verify cron expression is correct
	if _, err := cron.ParseStandard(val); err != nil {
		log.Fatalf("Invalid cron expression for %s: %v", key, err)
	}

	return val
}
