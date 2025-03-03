package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"userservices/internal/clients"
	"userservices/internal/models"

	"github.com/gorilla/mux"
)

type contextKey string

// Export the context keys
const (
	UserUUIDKey contextKey = "userUUID"
	UserRoleKey contextKey = "userRole"
)

// Add user info to context
func addUserToContext(ctx context.Context, uuid, role string) context.Context {
	ctx = context.WithValue(ctx, UserUUIDKey, uuid)
	return context.WithValue(ctx, UserRoleKey, role)
}

// Get user UUID from context
func GetUserUUIDFromContext(ctx context.Context) string {
	if uuid, ok := ctx.Value(UserUUIDKey).(string); ok {
		return uuid
	}
	return ""
}

// Get user role from context
func GetUserRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

func AuthMiddleware(authClient *clients.AuthClient) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Skip auth for login and register endpoints
            if r.URL.Path == "/login" || r.URL.Path == "/register" {
                next.ServeHTTP(w, r)
                return
            }

            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                sendAuthError(w, "No token provided", http.StatusUnauthorized)
                return
            }

            tokenParts := strings.Split(authHeader, " ")
            if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
                sendAuthError(w, "Invalid token format", http.StatusUnauthorized)
                return
            }

            // Validate token using AuthClient
            validationResp, err := authClient.ValidateToken(tokenParts[1])
            if err != nil {
                sendAuthError(w, "Error validating token", http.StatusInternalServerError)
                return
            }

            if validationResp.Body.IsValid != 1 {
                sendAuthError(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            // Add user info to request context
            r = r.WithContext(addUserToContext(r.Context(), validationResp.Body.UUID, validationResp.Body.Role))
            next.ServeHTTP(w, r)
        })
    }
}

// AdminOnly middleware checks if the user has admin role
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole := GetUserRoleFromContext(r.Context())
		if userRole != "admin" {
			sendAuthError(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func sendAuthError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.Response{
		Message:    message,
		StatusCode: statusCode,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
