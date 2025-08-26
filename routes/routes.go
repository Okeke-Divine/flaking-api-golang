package routes

import (
	"time"
	"github.com/Okeke-Divine/flaking-api/controllers"
	"github.com/Okeke-Divine/flaking-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize controllers
	userController := controllers.NewUserController()
	authController := controllers.NewAuthController()
	postController := controllers.NewPostController()
	rateLimiter := middleware.NewRateLimiter()
	go rateLimiter.CleanupOldVisits()

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes (no auth required)
		public := v1.Group("/auth")
		{
			public.POST("/register", rateLimiter.LimitBy(100, time.Minute), authController.Register) // 3 registrations per minute
			public.POST("/login", rateLimiter.LimitBy(100, time.Minute), authController.Login)       // 10 login attempts per minute
		}

		// Public user routes (read-only, no auth required)
		publicUsers := v1.Group("/users")
		{
			publicUsers.GET("", rateLimiter.LimitBy(300, time.Minute), userController.GetAllUsers)     // 30 requests per minute
			publicUsers.GET("/:id", rateLimiter.LimitBy(600, time.Minute), userController.GetUserByID) // 60 requests per minute
		}

		publicPosts := v1.Group("/posts")
		{
			publicPosts.GET("", rateLimiter.LimitBy(300, time.Minute), postController.GetAllPosts)     // NEW
			publicPosts.GET("/:id", rateLimiter.LimitBy(600, time.Minute), postController.GetPostByID) // NEW
		}

		// Protected routes (require JWT authentication)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User management
			user := protected.Group("/user")
			{
				user.GET("/profile", userController.GetCurrentUserProfile)          
				user.PUT("/profile", userController.UpdateUser)                 
				user.DELETE("/profile", userController.DeleteUser)                 
			}

			
			// Post management (authenticated users)
			posts := protected.Group("/posts")
			{
				posts.POST("", rateLimiter.LimitBy(30, time.Minute), postController.CreatePost)   
				posts.PUT("/:id", postController.UpdatePost)                                      
				posts.DELETE("/:id", postController.DeletePost)                            
			}

			// Admin routes
			admin := protected.Group("/admin")
			{
				admin.GET("/users", userController.GetAllUsers)
				admin.GET("/posts", postController.GetAllPosts)
			}
		}

		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "Server is running healthy!",
			})
		})
	}

	return router
}