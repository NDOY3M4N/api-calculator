package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

func Authentication(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := &wrapperWritter{w, http.StatusOK}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			logger.Info("Authorization header missing",
				slog.Int("statusCode", http.StatusUnauthorized),
				slog.String("remoteAddr", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("requestID", generateRequestID()),
			)

			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !isValidToken(token) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// If the token is valid, call the next handler
		next.ServeHTTP(wrapper, r)
	})
}

func isValidToken(token string) bool {
	// Implement your token verification logic here
	// For example, you could check the token against a database or use a JWT library
	return token == "valid_token"
}
