package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type wrapperWritter struct {
	http.ResponseWriter
	statusCode int
}

func generateRequestID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// Handle error
		return ""
	}
	return hex.EncodeToString(b)
}
