package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "2fa-go/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("❌ Error loading .env file")
    }

    // Construct the DSN (Data Source Name) using environment variables
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    log.Println("⌛ Connecting to database...")

    // Connect to the database
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("❌ Failed to connect to database:", err)
    }

    log.Println("✅ Successfully connected to PostgreSQL database")

    log.Println("🔄 Checking database migrations...")

        if !db.Migrator().HasTable(&models.User{}) {
            log.Println("🔄 Running database migrations...")
            err = db.AutoMigrate(&models.User{})
            if err != nil {
                log.Fatal("❌ Database migration failed:", err)
            }
            log.Println("👍 Database migrations completed successfully")
        } else {
            log.Println("✅ Database tables already exist - skipping migration")
        }

    DB = db

}