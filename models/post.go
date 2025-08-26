package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"not null;size:255" validate:"required,min=5,max=255"`
    Content   string         `json:"content" gorm:"type:text;not null" validate:"required,min=10"`
    UserID    uint           `json:"user_id" gorm:"not null;index" validate:"required"`
    User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"` // Relationship
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PostCreateRequest represents the data needed to create a post
type PostCreateRequest struct {
    Title   string `json:"title" binding:"required,min=5,max=255"`
    Content string `json:"content" binding:"required,min=10"`
}

// PostUpdateRequest represents the data needed to update a post
type PostUpdateRequest struct {
    Title   string `json:"title" binding:"omitempty,min=5,max=255"`
    Content string `json:"content" binding:"omitempty,min=10"`
}