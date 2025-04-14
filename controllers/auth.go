package controllers

import (
	"net/http"
	"2fa-go/config"
	"2fa-go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	OTPEnabled bool   `json:"otp_enabled"`
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if user.Email == "" || user.Password == "" || user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user.Password = string(hashedPassword)

	if result := config.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.Where("email = ?", credentials.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	responseUser := UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		OTPEnabled: user.OTPEnabled,
	}

	if user.OTPEnabled {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"otp":    true,
			"user":   responseUser,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("otp_verified", true)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   responseUser,
	})
}

func GenerateOTP(c *gin.Context) {
	var input struct {
		UserID uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.First(&user, input.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "2FA Demo App",
		AccountName: user.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating OTP"})
		return
	}

	config.DB.Model(&user).Updates(models.User{OTPSecret: key.Secret()})

	c.JSON(http.StatusOK, gin.H{
		"otpauth_url": key.URL(),
		"base32":      key.Secret(),
	})
}

func VerifyOTP(c *gin.Context) {
	var input struct {
		Token  string `json:"token"`
		UserID uint   `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.First(&user, input.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	valid := totp.Validate(input.Token, user.OTPSecret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	config.DB.Model(&user).Update("otp_enabled", true)

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("otp_verified", true)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	responseUser := UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		OTPEnabled: true,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "2FA enabled successfully",
		"user":    responseUser,
	})
}

func ValidateOTP(c *gin.Context) {
	var input struct {
		Token  string `json:"token"`
		UserID uint   `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.First(&user, input.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	valid := totp.Validate(input.Token, user.OTPSecret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("otp_verified", true)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	responseUser := UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		OTPEnabled: user.OTPEnabled,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   responseUser,
	})
}

func DisableOTP(c *gin.Context) {
	session := sessions.Default(c)
	sessionUserID := session.Get("user_id")
	if sessionUserID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var input struct {
		Password string `json:"password"`
		UserID   uint   `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if sessionUserID != input.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to perform this action"})
		return
	}

	var user models.User
	if result := config.DB.First(&user, input.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	updates := map[string]interface{}{
		"otp_enabled": false,
		"otp_secret":  "",
	}

	if result := config.DB.Model(&user).Updates(updates); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	session.Set("otp_verified", false)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "2FA disabled successfully",
		"user": UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			OTPEnabled: false,
		},
	})
}

func Profile(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var user models.User
	if result := config.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			OTPEnabled: user.OTPEnabled,
		},
	})
}