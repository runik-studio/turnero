package auth

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"ServiceBookingApp/internal/domain"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// AuthService defines the interface for authentication
type AuthService interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

// FirebaseAuthService implements AuthService using Firebase
type FirebaseAuthService struct {
	Client *auth.Client
}

func (s *FirebaseAuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return s.Client.VerifyIDToken(ctx, idToken)
}

// MockAuthService implements AuthService for testing
type MockAuthService struct{}

func (m *MockAuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	// Return a valid mock token
	return &auth.Token{
		UID: "test-user-id",
		Claims: map[string]interface{}{
			"email": "test@example.com",
			"name":  "Test User",
		},
	}, nil
}

// AuthMiddleware verifies the Firebase ID token
func AuthMiddleware(service AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			return
		}

		// Check for MOCK_AUTH
		if os.Getenv("MOCK_AUTH") == "true" && tokenString == "mock-token" {
			c.Set("user", &auth.Token{
				UID: "mock-user-id",
				Claims: map[string]interface{}{
					"email": "mock@example.com",
					"role":  "admin",
				},
			})
			c.Next()
			return
		}

		token, err := service.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			return
		}

		// Store user info in context
		c.Set("user", token)
		c.Next()
	}
}

func UserActiveMiddleware(repo domain.UsersRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userTokenInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		userToken, ok := userTokenInterface.(*auth.Token)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user token type"})
			return
		}
		uid := userToken.UID

		user, err := repo.Get(c.Request.Context(), uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User profile not found or access denied"})
			return
		}

		if user.IsActive != nil && !*user.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Account is blocked"})
			return
		}

		c.Set("user_data", user)
		c.Next()
	}
}
