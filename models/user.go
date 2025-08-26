package models

import (
    "time"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
    Password  string         `json:"-" gorm:"not null" validate:"required,min=6"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    return u.HashPassword()
}

// HashPassword encrypts the user password
func (u *User) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

// CheckPassword compares plain password with hashed password
func (u *User) CheckPassword(plainPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
}

// UserCreateRequest represents the data needed to create a user
type UserCreateRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateRequest represents the data needed to update a user
type UserUpdateRequest struct {
    Name  string `json:"name" binding:"omitempty,min=2,max=100"`
    Email string `json:"email" binding:"omitempty,email"`
}