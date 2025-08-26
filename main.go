package main

import (
    "log"
    "strconv"
    "github.com/Okeke-Divine/flaking-api/config"
    "github.com/Okeke-Divine/flaking-api/database"
    "github.com/Okeke-Divine/flaking-api/routes"
    "github.com/joho/godotenv"
)

func main() {
    // Load configuration
    godotenv.Load()
    config.LoadConfig()

    // Initialize database
    err := database.Connect()
    if err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }

    // Run auto migrations
    err = database.AutoMigrate()
    if err != nil {
        log.Fatalf("Auto migration failed: %v", err)
    }

    // Set up routes
    router := routes.SetupRouter()
    router.SetTrustedProxies([]string{"127.0.0.1"})

    // Start server
    log.Printf("Server starting on port %d", config.AppConfig.Port)
    router.Run(":" + strconv.Itoa(config.AppConfig.Port))
}