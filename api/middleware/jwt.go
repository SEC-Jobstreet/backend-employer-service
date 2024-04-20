package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// VerifyToken verifies a JWT token using Cognito JWKS
func VerifyToken(tokenString string, jwksUrl string) error {
	// Fetch JWKS from Cognito
	ar := jwk.NewAutoRefresh(context.Background())
	ar.Configure(jwksUrl, jwk.WithMinRefreshInterval(15*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	set, err := ar.Refresh(ctx, jwksUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	// Parse and verify the token
	parsedToken, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(set),
		jwt.WithValidate(true),
	)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	// Additional validation can be performed here (e.g., checking 'iss', 'aud', etc.)

	fmt.Printf("Token claims: %#v\n", parsedToken)
	return nil
}

// CognitoAuthMiddleware creates a Gin middleware for JWT validation using Cognito
func CognitoAuthMiddleware(jwksUrl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Call VerifyToken here
		err := VerifyToken(bearerToken, jwksUrl)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// If the token is valid, proceed to the next handler
		c.Next()
	}
}
