package routes

import (
	"2fa-go/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/otp/generate", controllers.GenerateOTP)
			auth.POST("/otp/verify", controllers.VerifyOTP)
			auth.POST("/otp/validate", controllers.ValidateOTP)
			auth.POST("/otp/disable", controllers.DisableOTP)
		}

		// Protected routes
		api.GET("/profile", controllers.Profile)
	}
}