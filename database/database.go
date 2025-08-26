package database

import (
    "fmt"
    "log"
    "github.com/Okeke-Divine/flaking-api/config"
    "github.com/Okeke-Divine/flaking-api/models"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
    // Step 1: Connect without database to create it if needed
    dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
        config.AppConfig.DBUser,
        config.AppConfig.DBPassword,
        config.AppConfig.DBHost,
        config.AppConfig.DBPort)

    tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to MySQL server: %w", err)
    }

    // Step 2: Create database if it doesn't exist
    createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.AppConfig.DBName)
    result := tempDB.Exec(createDBSQL)
    if result.Error != nil {
        return fmt.Errorf("failed to create database: %w", result.Error)
    }
    log.Printf("Database '%s' ensured", config.AppConfig.DBName)

    // Close temporary connection
    sqlDB, _ := tempDB.DB()
    sqlDB.Close()

    // Step 3: Connect to the actual database
    dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        config.AppConfig.DBUser,
        config.AppConfig.DBPassword,
        config.AppConfig.DBHost,
        config.AppConfig.DBPort,
        config.AppConfig.DBName)

    DB, err = gorm.Open(mysql.Open(dsnWithDB), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    log.Println("Database connection established")
    return nil
}

func AutoMigrate() error {
    // Skip auto migration in production
    if config.AppConfig.Environment == "production" {
        log.Println("Skipping auto migrate in production - use SQL migrations")
        return nil
    }
    
    modelsToMigrate := models.RegisterModels()
    
    err := DB.AutoMigrate(modelsToMigrate...)
    if err != nil {
        return fmt.Errorf("failed to auto migrate: %w", err)
    }
    
    log.Printf("Auto migration completed for %d models", len(modelsToMigrate))
    return nil
}