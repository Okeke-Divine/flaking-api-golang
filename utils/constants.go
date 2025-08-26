package utils

const (
    APIVersion      = "1.0"
    MaxPageSize     = 100
    DefaultPageSize = 10
)

var (
    AppName    = "Flaking API"
    AppVersion = "1.0.0"
    JWTSecretKey string
)

const (
    ErrCodeValidation = "VALIDATION_ERROR"
    ErrCodeNotFound   = "NOT_FOUND"
    ErrCodeServer     = "INTERNAL_ERROR"
)