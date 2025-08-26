package config

import (
    "log"
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    Port        int
    Environment string
    DBHost      string
    DBPort      int
    DBName      string
    DBUser      string
    DBPassword  string
    JWTSecret   string
}

var AppConfig Config

func LoadConfig() {
    if os.Getenv("ENVIRONMENT") != "production" {
        err := godotenv.Load()
        if err != nil {
            log.Println("Warning: No .env file found")
        }
    }
    
    AppConfig = Config{
        Port:        getEnvAsInt("PORT", 3000),
        Environment: getEnv("ENVIRONMENT", "development"),
        DBHost:      getEnv("DB_HOST", "localhost"),
        DBPort:      getEnvAsInt("DB_PORT", 3306), // MySQL default port
        DBName:      getEnv("DB_NAME", "flaking_db"),
        DBUser:      getEnv("DB_USER", "root"),
        DBPassword:  getEnv("DB_PASSWORD", ""),
        JWTSecret:   getEnv("JWT_SECRET", "fallback-secret-change-in-production"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}