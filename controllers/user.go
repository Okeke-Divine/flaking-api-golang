package controllers

import (
    "net/http"
    "strconv"

    "github.com/Okeke-Divine/flaking-api/database"
    "github.com/Okeke-Divine/flaking-api/models"
    "github.com/Okeke-Divine/flaking-api/utils"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "fmt"
)

type UserController struct {
    DB *gorm.DB
}

func NewUserController() *UserController {
    return &UserController{DB: database.DB}
}

// Add this helper method to get current user ID from context
func getCurrentUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("not authenticated")
	}
	
	return userID.(uint), nil
}


// GetAllUsers handles GET /users - Get all users with pagination
func (uc *UserController) GetAllUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    
    offset, limit := utils.Pagination(page, pageSize)
    
    var users []models.User
    var total int64

    // Get total count
    uc.DB.Model(&models.User{}).Count(&total)
    
    // Get paginated users
    result := uc.DB.Offset(offset).Limit(limit).Find(&users)
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
        return
    }

    utils.SuccessResponse(c, "Users retrieved successfully", gin.H{
        "users": users,
        "pagination": gin.H{
            "page":       page,
            "pageSize":   limit,
            "total":      total,
            "totalPages": (int(total) + limit - 1) / limit,
        },
    })
}

// GetUserByID handles GET /users/:id - Get user by ID
func (uc *UserController) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    
    var user models.User
    result := uc.DB.First(&user, id)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            utils.ErrorResponse(c, http.StatusNotFound, "User not found")
        } else {
            utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user")
        }
        return
    }

    utils.SuccessResponse(c, "User retrieved successfully", user)
}

// CreateUser handles POST /users - Create new user
func (uc *UserController) CreateUser(c *gin.Context) {
    var request models.UserCreateRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
        return
    }

    // Check if email already exists
    var existingUser models.User
    if uc.DB.Where("email = ?", request.Email).First(&existingUser).Error == nil {
        utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
        return
    }

    // Create user
    user := models.User{
        Name:     request.Name,
        Email:    request.Email,
        Password: request.Password,
    }

    result := uc.DB.Create(&user)
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
        return
    }

    // Don't return password in response
    user.Password = ""
    
    utils.SuccessResponse(c, "User created successfully", user)
}

// UpdateUser handles PUT /users/:id - Update user (now requires ownership)
func (uc *UserController) UpdateUser(c *gin.Context) {
	currentUserID, err := getCurrentUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	var request models.UserUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	var user models.User
	result := uc.DB.First(&user, currentUserID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	// Update fields if provided
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Email != "" {
		user.Email = request.Email
	}

	result = uc.DB.Save(&user)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.SuccessResponse(c, "User updated successfully", user)
}


// DeleteUser handles DELETE /users/:id - Delete user (requires ownership)
func (uc *UserController) DeleteUser(c *gin.Context) {
	currentUserID, err := getCurrentUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	id := c.Param("id")
	
	// Convert param ID to uint
	paramID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Check if user is deleting their own profile
	if uint(paramID) != currentUserID {
		utils.ErrorResponse(c, http.StatusForbidden, "Can only delete your own profile")
		return
	}

	result := uc.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	if result.RowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, "User deleted successfully", nil)
}

// Add this new method for getting current user's profile
func (uc *UserController) GetCurrentUserProfile(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	var user models.User
	result := uc.DB.First(&user, userID)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Don't return password in response
	user.Password = ""

	utils.SuccessResponse(c, "User profile retrieved successfully", user)
}

