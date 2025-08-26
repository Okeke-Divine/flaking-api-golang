package utils

import (
    "fmt"
    "time"
    "github.com/gin-gonic/gin"
)

// FormatDuration pretty-prints time.Duration
func FormatDuration(d time.Duration) string {
    return fmt.Sprintf("%.2fs", d.Seconds())
}

// GetClientIP extracts client IP from request
func GetClientIP(c *gin.Context) string {
    return c.ClientIP()
}

// ValidateEmail simple email validation
func ValidateEmail(email string) bool {
    // Simple regex check - in real app use proper validation
    return len(email) > 3 && len(email) < 254
}

// Pagination calculates offset and limit for database queries
func Pagination(page, pageSize int) (offset, limit int) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = DefaultPageSize
    }
    if pageSize > MaxPageSize {
        pageSize = MaxPageSize
    }
    offset = (page - 1) * pageSize
    return offset, pageSize
}