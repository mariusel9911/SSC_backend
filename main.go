package main

import (
    "2fa-go/config"
    "2fa-go/routes"
    "net/http"
    "time"
    "log"
    "os"

    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    //.env
    err := godotenv.Load()
    if err != nil {
        log.Fatal("❌ Error loading .env file")
    }

    config.ConnectDB()

    r := gin.Default()

    // Session configuration
    store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")));
    store.Options(sessions.Options{
        Path:     "/",
        MaxAge:   86400 * 7, // 86400 secunde = 1 zi
        Secure:   false,    // false = HTTP -> in development
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode, // frontend/backend -> different ports
    })
    r.Use(sessions.Sessions("2fa_session", store))

    // CORS configuration
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Set-Cookie"},
        ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    routes.SetupRoutes(r)

    r.Run(":8000")
}