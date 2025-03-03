package models

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"-"` // "-" means this won't be included in JSON
    Role     string `json:"role"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     string `json:"role,omitempty"` // Optional role field
}

type Response struct {
    Message    string      `json:"message"`
    Body       interface{} `json:"body,omitempty"`
    StatusCode int         `json:"status_code"`
}

type UserResponse struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    UUID  string `json:"uuid"`
    Role  string `json:"role"`
}

type LoginResponse struct {
    User  UserResponse `json:"user"`
    Token string      `json:"token"`
}

type ValidationBody struct {
    IsValid int    `json:"isValid"`
    UUID    string `json:"uuid"`
    Role    string `json:"role"`
}

type ValidationResponse struct {
    Message    string         `json:"message"`
    Body       ValidationBody `json:"body"`
    StatusCode int           `json:"status_code"`
}