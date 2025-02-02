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

type contextKey string

const (
	UserIDContextKey   contextKey = "userID"
	UserNameContextKey contextKey = "userName"
)

type Middleware struct {
	cfg config.Config
	log *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) *Middleware {
	return &Middleware{cfg: cfg, log: log}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accesstoken")
		if err != nil {
			m.log.Warn("No access token cookie", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		secretKey := []byte(m.cfg.AccessSecret)

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

		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			m.log.Warn("Token expired", zap.Time("expiresAt", claims.ExpiresAt.Time))
			http.Error(w, "unauthorized: token expired", http.StatusUnauthorized)
			return
		}

		m.log.Info("Token is valid", zap.Int64("userID", claims.UserID), zap.String("username", claims.Username))

		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, UserNameContextKey, claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
