package models

// RegisterModels returns all models for auto migration
func RegisterModels() []interface{} {
    return []interface{}{
        &User{},
        &Post{},
    }
}