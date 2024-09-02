package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	IsDevelopment bool
	Server        ServerConfig
	Database      DbConfig
}

type ServerConfig struct {
	Port         string
	JwtSecret    string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	LogPath      string
	LogLevel     string
	UploadDir    string

	JwtAccessTokenExpireMinute time.Duration
	JwtRefreshTokenExpireHour  time.Duration
}

type DbConfig struct {
	DbUsername  string
	DbPassword  string
	DbHost      string
	DbPort      string
	DbName      string
	MaxPoolSize string
	MaxIdleConn string
	MaxLifetime string
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize configuration struct
	config := &Config{
		IsDevelopment: false,
		Server: ServerConfig{
			Port:                       getEnv("SERVER_PORT", "3003"),
			JwtSecret:                  getEnv("SERVER_SECRET", ""),
			JwtAccessTokenExpireMinute: getEnvDuration("JWT_ACCESS_TOKEN_EXPIRE_MINUTE", "15m"),
			JwtRefreshTokenExpireHour:  getEnvDuration("JWT_REFRESH_TOKEN_EXPIRE_HOUR", "24h"),
		},
		Database: DbConfig{
			DbUsername:  getEnv("DB_USERNAME", "username"),
			DbPassword:  getEnv("DB_PASSWORD", "password"),
			DbHost:      getEnv("DB_HOST", "localhost"),
			DbPort:      getEnv("DB_PORT", "3306"),
			DbName:      getEnv("DB_NAME", "dbname"),
			MaxIdleConn: getEnv("MAX_IDLE_CONN", "1"),
			MaxPoolSize: getEnv("MAX_POOL_SIZE", "5"),
			MaxLifetime: getEnv("MAX_LIFE_TIME", "1800"),
		},
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvDuration(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		return time.Duration(0) // Hatalı format durumunda varsayılan olarak 0 döner
	}
	return duration
}
