package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPasswd               string
	DBAddr                 string
	DBName                 string
	FrontendURL            string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPasswd:   getEnv("DB_PASSWD", "root"),
		DBAddr: fmt.Sprintf(
			"%s:%s", getEnv("DB_HOST", "127.0.0.1:3306"), getEnv("DB_PORT", "3306"),
		),
		DBName:                 getEnv("DB_NAME", "ecom"),
		FrontendURL:            getEnv("FRONT_URL", "http://localhost:3000"),
		JWTSecret:              getEnv("JWT_SECRET", "not-secret-anymore?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
