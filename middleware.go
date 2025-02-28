package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/NDOY3M4N/api-calculator/ratelimit"
	"github.com/NDOY3M4N/api-calculator/repository"
)

type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type wrapperWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapperWritter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func IsAuthenticated(repo *repository.Repository) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeError(
					w,
					r,
					http.StatusUnauthorized,
					fmt.Errorf("missing authorization header"),
				)
				return
			}

			token, err := ValidateToken(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				writeError(
					w,
					r,
					http.StatusForbidden,
					fmt.Errorf("error validating token: %s", err),
				)
				return
			}

			if !token.Valid {
				writeError(w, r, http.StatusForbidden, fmt.Errorf("token invalid"))
				return
			}

			var userIDStr string
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userIDStr = claims["userID"].(string)
			} else {
				writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to get claims from token"))
				return
			}

			userID, _ := strconv.Atoi(userIDStr)
			if _, err = repo.FindUserById(userID); err != nil {
				writeError(w, r, http.StatusUnauthorized, fmt.Errorf("permission denied"))
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func RateLimit(tb *ratelimit.TokenBucket) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Context().Value(requestIDKey).(string)
			tb.Consume()

			remaining := len(tb.Tokens)
			reset := time.Now().Add(time.Second).Unix()
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", bucketSize))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", reset))

			if remaining == 0 {
				logger.Warn("Rate limit exceeded. Please wait before making more requests.",
					slog.Int("statusCode", http.StatusTooManyRequests),
					slog.String("remoteAddr", r.RemoteAddr),
					slog.Group("request",
						slog.String("id", reqID),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					),
				)

				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func AddRequestId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, err := generateRequestID()
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		r = r.WithContext(ctx)

		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r)
	}
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(requestIDKey).(string)

		start := time.Now()
		wrappedWritter := &wrapperWritter{w, http.StatusOK}

		logger.Info("Log request",
			slog.Int("statusCode", wrappedWritter.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.Group("request",
				slog.String("id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			),
		)

		next.ServeHTTP(wrappedWritter, r)
	}
}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

func generateRequestID() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Error generating request ID: %s", err)
	}

	return "req_" + hex.EncodeToString(b), nil
}
