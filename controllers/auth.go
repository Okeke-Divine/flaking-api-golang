package controllers

import (
	"net/http"
	"github.com/Okeke-Divine/flaking-api/database"
	"github.com/Okeke-Divine/flaking-api/models"
	"github.com/Okeke-Divine/flaking-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"fmt"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController() *AuthController {
	return &AuthController{DB: database.DB}
}

// LoginRequest represents login input
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var request models.UserCreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	// Check if email already exists
	var existingUser models.User
	if ac.DB.Where("email = ?", request.Email).First(&existingUser).Error == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	// Create user
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password, // Will be hashed by BeforeCreate hook
	}

	fmt.Printf("Creating user: %+v\n", user)

	result := ac.DB.Create(&user)
	if result.Error != nil {
		fmt.Printf("Error: %v \n", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Don't return password in response
	user.Password = ""

	utils.SuccessResponse(c, "User registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	// Find user by email
	var user models.User
	result := ac.DB.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	// Check password
	if err := user.CheckPassword(request.Password); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Don't return password in response
	user.Password = ""

	utils.SuccessResponse(c, "Login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

// GetCurrentUser returns current authenticated user
func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var user models.User
	result := ac.DB.First(&user, userID)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Don't return password in response
	user.Password = ""

	utils.SuccessResponse(c, "Current user retrieved", user)
}