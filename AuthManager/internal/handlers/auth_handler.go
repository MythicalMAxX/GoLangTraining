package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Action string `json:"action"`
	Token  string `json:"token"`
}

type AuthResponse struct {
	Token     string `json:"token,omitempty"`
	IsValid   int    `json:"isValid,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Role      string `json:"role,omitempty"`
	Action    string `json:"action"`
	RequestID string `json:"request_id"`
}

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("START: GenerateToken request from %s", r.RemoteAddr)

	// Check if SECRET_KEY is set
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Printf("ERROR: SECRET_KEY environment variable is not set")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: SECRET_KEY is properly configured")

	// Decode request
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Successfully decoded request - UserID: %s, Role: %s, Action: %s",
		req.UserID, req.Role, req.Action)

	// Validate request fields
	if req.UserID == "" || req.Role == "" {
		log.Printf("ERROR: Invalid request - UserID or Role is empty")
		http.Error(w, "UserID and Role are required", http.StatusBadRequest)
		return
	}

	// Generate token
	token := h.generateJWT(req.UserID, req.Role)
	if token == "" {
		log.Printf("ERROR: Failed to generate JWT token")
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: Successfully generated JWT token")

	// Prepare response
	response := AuthResponse{
		Token:     token,
		Action:    "generate",
		RequestID: uuid.New().String(),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ERROR: Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("SUCCESS: GenerateToken completed for RequestID: %s", response.RequestID)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("START: ValidateToken request from %s", r.RemoteAddr)

	// Check if SECRET_KEY is set
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Printf("ERROR: SECRET_KEY environment variable is not set")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: SECRET_KEY is properly configured")

	// Decode request
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Successfully decoded request - Token: %s, Action: %s",
		req.Token, req.Action)

	// Validate token
	claims, valid := h.validateJWT(req.Token)
	log.Printf("INFO: Token validation result - Valid: %v", valid)

	// Handle claims
	var userID, role string
	if valid && claims != nil {
		if uid, ok := claims["user_id"].(string); ok {
			userID = uid
			log.Printf("INFO: Extracted UserID from token: %s", userID)
		}
		if r, ok := claims["role"].(string); ok {
			role = r
			log.Printf("INFO: Extracted Role from token: %s", role)
		}
	}

	// Prepare response
	response := AuthResponse{
		IsValid:   boolToInt(valid),
		UUID:      userID,
		Role:      role,
		Action:    "validate",
		RequestID: uuid.New().String(),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ERROR: Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("SUCCESS: ValidateToken completed for RequestID: %s - Valid: %v, UserID: %s, Role: %s",
		response.RequestID, valid, userID, role)
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (h *AuthHandler) generateJWT(userID, role string) string {
	log.Printf("START: Generating JWT for UserID: %s, Role: %s", userID, role)

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Printf("ERROR: SECRET_KEY is empty")
		return ""
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	log.Printf("INFO: Created claims with expiry at: %v", time.Unix(claims["exp"].(int64), 0))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("ERROR: Failed to sign token: %v", err)
		return ""
	}

	log.Printf("SUCCESS: JWT token generated successfully")
	return tokenString
}

func (h *AuthHandler) validateJWT(tokenString string) (jwt.MapClaims, bool) {
	log.Printf("START: Validating JWT token")

	if tokenString == "" {
		log.Printf("ERROR: Token string is empty")
		return make(jwt.MapClaims), false
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Printf("ERROR: SECRET_KEY is empty")
		return make(jwt.MapClaims), false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			log.Printf("ERROR: %s", msg)
			return nil, fmt.Errorf(msg)
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("ERROR: Token parsing failed: %v", err)
		return make(jwt.MapClaims), false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("SUCCESS: Token is valid - Claims: %v", claims)
		return claims, true
	}

	log.Printf("ERROR: Token is invalid or claims could not be extracted")
	return make(jwt.MapClaims), false
}
