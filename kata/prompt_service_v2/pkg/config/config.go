// pkg/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	JwtSecret      string
	DatabasePath   string
	AuthzModelPath string

	RedisHost     string
	RedisPort     int16
	RedisPassword string
}

func LoadConfig() (*Config, error) {
	portStr := os.Getenv("REDIS_PORT")
	if portStr == "" {
		portStr = "8080" // default
	}

	redisPort, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil || redisPort < 1 || redisPort > 65535 {
		return nil, fmt.Errorf("invalid PORT value: %s", portStr)
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "prompt.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, &MissingEnvError{"JWT_SECRET"}
	}

	authzModelPath := os.Getenv("AUTHZ_MODEL_PATH")
	if authzModelPath == "" {
		return nil, &MissingEnvError{"AUTHZ_MODEL_PATH"}
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	return &Config{
		JwtSecret:      jwtSecret,
		DatabasePath:   dbPath,
		AuthzModelPath: authzModelPath,
		RedisHost:      redisHost,
		RedisPort:      int16(redisPort),
		RedisPassword:  redisPassword,
	}, nil
}

type MissingEnvError struct {
	VarName string
}

func (e *MissingEnvError) Error() string {
	return "missing environment variable: " + e.VarName
}
