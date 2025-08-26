package controllers

import (
    "net/http"
    "strconv"

    "github.com/Okeke-Divine/flaking-api/database"
    "github.com/Okeke-Divine/flaking-api/models"
    "github.com/Okeke-Divine/flaking-api/utils"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type PostController struct {
    DB *gorm.DB
}

func NewPostController() *PostController {
    return &PostController{DB: database.DB}
}

// GetAllPosts handles GET /posts - Get all posts with pagination
func (pc *PostController) GetAllPosts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    
    offset, limit := utils.Pagination(page, pageSize)
    
    var posts []models.Post
    var total int64

    // Get total count
    pc.DB.Model(&models.Post{}).Count(&total)
    
    // Get paginated posts with user information
    result := pc.DB.Offset(offset).Limit(limit).Preload("User").Find(&posts)
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch posts")
        return
    }

    utils.SuccessResponse(c, "Posts retrieved successfully", gin.H{
        "posts": posts,
        "pagination": gin.H{
            "page":       page,
            "pageSize":   limit,
            "total":      total,
            "totalPages": (int(total) + limit - 1) / limit,
        },
    })
}

// GetPostByID handles GET /posts/:id - Get post by ID
func (pc *PostController) GetPostByID(c *gin.Context) {
    id := c.Param("id")
    
    var post models.Post
    result := pc.DB.Preload("User").First(&post, id)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
        } else {
            utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch post")
        }
        return
    }

    utils.SuccessResponse(c, "Post retrieved successfully", post)
}

// GetPostsByCurrentUser handles GET /user/posts - Get current user's posts
func (pc *PostController) GetPostsByCurrentUser(c *gin.Context) {
    userID, err := getCurrentUserID(c)
    if err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    
    offset, limit := utils.Pagination(page, pageSize)
    
    var posts []models.Post
    var total int64

    // Get total count for current user
    pc.DB.Model(&models.Post{}).Where("user_id = ?", userID).Count(&total)
    
    // Get paginated posts for current user
    result := pc.DB.Offset(offset).Limit(limit).
        Where("user_id = ?", userID).
        Preload("User").
        Find(&posts)
        
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch your posts")
        return
    }

    utils.SuccessResponse(c, "Your posts retrieved successfully", gin.H{
        "posts": posts,
        "pagination": gin.H{
            "page":       page,
            "pageSize":   limit,
            "total":      total,
            "totalPages": (int(total) + limit - 1) / limit,
        },
    })
}

// CreatePost handles POST /posts - Create new post (authenticated users only)
func (pc *PostController) CreatePost(c *gin.Context) {
    userID, err := getCurrentUserID(c)
    if err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
        return
    }

    var request models.PostCreateRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
        return
    }

    // Create post
    post := models.Post{
        Title:   request.Title,
        Content: request.Content,
        UserID:  userID, // Set the foreign key
    }

    result := pc.DB.Create(&post)
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
        return
    }

    // Preload user data for response
    pc.DB.Preload("User").First(&post, post.ID)

    utils.SuccessResponse(c, "Post created successfully", post)
}

// UpdatePost handles PUT /posts/:id - Update post (only by owner)
func (pc *PostController) UpdatePost(c *gin.Context) {
    userID, err := getCurrentUserID(c)
    if err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
        return
    }

    id := c.Param("id")
    
    var request models.PostUpdateRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
        return
    }

    // First, verify the post exists and belongs to the current user
    var post models.Post
    result := pc.DB.Where("id = ? AND user_id = ?", id, userID).First(&post)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            utils.ErrorResponse(c, http.StatusNotFound, "Post not found or you don't have permission to edit it")
        } else {
            utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch post")
        }
        return
    }

    // Update fields if provided
    if request.Title != "" {
        post.Title = request.Title
    }
    if request.Content != "" {
        post.Content = request.Content
    }

    result = pc.DB.Save(&post)
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post")
        return
    }

    // Preload user data for response
    pc.DB.Preload("User").First(&post, post.ID)

    utils.SuccessResponse(c, "Post updated successfully", post)
}

// DeletePost handles DELETE /posts/:id - Delete post (only by owner)
func (pc *PostController) DeletePost(c *gin.Context) {
    userID, err := getCurrentUserID(c)
    if err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
        return
    }

    id := c.Param("id")
    
    // Delete only if post belongs to current user
    result := pc.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Post{})
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post")
        return
    }

    if result.RowsAffected == 0 {
        utils.ErrorResponse(c, http.StatusNotFound, "Post not found or you don't have permission to delete it")
        return
    }

    utils.SuccessResponse(c, "Post deleted successfully", nil)
}

// GetPostsByUserID handles GET /users/:id/posts - Get posts by specific user
func (pc *PostController) GetPostsByUserID(c *gin.Context) {
    userID := c.Param("id")
    
    // Verify user exists
    var user models.User
    if err := pc.DB.First(&user, userID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.ErrorResponse(c, http.StatusNotFound, "User not found")
        } else {
            utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user")
        }
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    
    offset, limit := utils.Pagination(page, pageSize)
    
    var posts []models.Post
    var total int64

    // Get total count for this user
    pc.DB.Model(&models.Post{}).Where("user_id = ?", userID).Count(&total)
    
    // Get paginated posts for this user
    result := pc.DB.Offset(offset).Limit(limit).
        Where("user_id = ?", userID).
        Preload("User").
        Find(&posts)
        
    if result.Error != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user's posts")
        return
    }

    utils.SuccessResponse(c, "User posts retrieved successfully", gin.H{
        "posts": posts,
        "user": gin.H{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
        "pagination": gin.H{
            "page":       page,
            "pageSize":   limit,
            "total":      total,
            "totalPages": (int(total) + limit - 1) / limit,
        },
    })
}