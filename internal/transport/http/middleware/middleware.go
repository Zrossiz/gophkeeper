// Package middleware provides HTTP middleware utilities, including authentication handling.
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Zrossiz/gophkeeper/internal/config"
	"github.com/Zrossiz/gophkeeper/internal/utils"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// contextKey represents a custom type for storing values in the request context.
type contextKey string

const (
	// UserIDContextKey is the context key for storing the user ID.
	UserIDContextKey contextKey = "userID"

	// UserNameContextKey is the context key for storing the username.
	UserNameContextKey contextKey = "userName"
)

// Middleware provides middleware functions for handling authentication and request processing.
type Middleware struct {
	cfg config.Config // Application configuration settings.
	log *zap.Logger   // Logger instance for logging events and errors.
}

// New creates a new Middleware instance.
//
// Parameters:
//   - cfg config.Config: The application configuration.
//   - log *zap.Logger: Logger for structured logging.
//
// Returns:
//   - *Middleware: A pointer to the initialized Middleware struct.
func New(cfg config.Config, log *zap.Logger) *Middleware {
	return &Middleware{cfg: cfg, log: log}
}

// Auth is an HTTP middleware that validates JWT tokens from request cookies.
//
// It extracts the "accesstoken" cookie, parses the JWT token, and validates its signature
// and expiration time. If the token is valid, the user ID and username are stored in the request context.
//
// Parameters:
//   - next http.Handler: The next handler to call after authentication.
//
// Returns:
//   - http.Handler: A handler that performs authentication before calling the next handler.
func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the access token from cookies.
		cookie, err := r.Cookie("accesstoken")
		if err != nil {
			m.log.Warn("No access token cookie", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		secretKey := []byte(m.cfg.AccessSecret)

		// Parse and validate the JWT token.
		claims := &utils.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			m.log.Warn("Token parsing failed", zap.Error(err))
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			m.log.Warn("Invalid token", zap.String("token", token.Raw))
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// Check if the token has expired.
		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			m.log.Warn("Token expired", zap.Time("expiresAt", claims.ExpiresAt.Time))
			http.Error(w, "unauthorized: token expired", http.StatusUnauthorized)
			return
		}

		m.log.Info("Token is valid", zap.Int64("userID", claims.UserID), zap.String("username", claims.Username))

		// Store user details in request context.
		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, UserNameContextKey, claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
