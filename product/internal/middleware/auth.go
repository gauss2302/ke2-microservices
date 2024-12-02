package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	pb "github.com/gauss2302/testcommm/product/proto/auth"
)

type AuthMiddleware struct {
	authClient pb.AuthServiceClient
}

func NewAuthMiddleware(authClient pb.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// product/internal/middleware/auth.go
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		log.Printf("Got auth header: %s", authHeader)

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Invalid auth header format")
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		log.Printf("Trying to verify token: %s", token)

		// Verify token with auth service
		resp, err := m.authClient.VerifyToken(r.Context(), &pb.VerifyTokenRequest{
			Token: token,
		})
		if err != nil {
			log.Printf("Error verifying token with auth service: %v", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Token verified successfully, user ID: %d", resp.UserId)
		ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
