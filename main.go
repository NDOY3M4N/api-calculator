package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/NDOY3M4N/api-calculator/docs"
)

type Payload struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

type APIError struct {
	Error string
}

type APISuccess struct {
	Result int
}

// addHandler Foo
//
// @summary Add two numbers
// @description Add two numbers together
// @tags Math
// @accept json
// @produce json
// @param payload body Payload true "Numbers needed for the operation"
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /add [post]
func addHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := parseJSON(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 + payload.Number2})
}

// substractHandler Foo
//
// @summary Substract two numbers
// @description Substract two numbers together
// @tags Math
// @accept json
// @produce json
// @param payload body Payload true "Numbers needed for the operation"
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /substract [post]
func substractHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := parseJSON(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 - payload.Number2})
}

// multiplyHandler Foo
//
// @summary Multiply two numbers
// @description Multiply two numbers together
// @tags Math
// @accept json
// @produce json
// @param payload body Payload true "Numbers needed for the operation"
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /multiply [post]
func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := parseJSON(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 * payload.Number2})
}

// divideHandler Foo
//
// @summary Divide two numbers
// @description Divide two numbers together
// @tags Math
// @accept json
// @produce json
// @param payload body Payload true "Numbers needed for the operation"
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /divide [post]
func divideHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := parseJSON(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if payload.Number2 == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Division by zero"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 / payload.Number2})
}

// @title         Calculator API
// @version       1.0
// @description   This is a sample server for Calculator API

// @contact.name  Abdoulaye NDOYE
// @contact.url   https://github.com/NDOY3M4N
// @contact.email pa.ndoye@outlook.com

// @license.name  MIT
// @license.url   https://github.com/NDOY3M4N/calculator-api/blob/main/LICENSE

// @host          localhost:3000
// @BasePath      /api/v1
func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /add", addHandler)
	router.HandleFunc("POST /substract", substractHandler)
	router.HandleFunc("POST /multiply", multiplyHandler)
	router.HandleFunc("POST /divide", divideHandler)

	// Create a separate handler for Swagger
	swaggerHandler := http.StripPrefix("/swagger/", http.HandlerFunc(
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(httpSwagger.HideModel),
		),
	))

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// Combine the v1 handler and the Swagger handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/" || strings.HasPrefix(r.URL.Path, "/swagger/") {
			swaggerHandler.ServeHTTP(w, r)
		} else {
			v1.ServeHTTP(w, r)
		}
	})

	server := http.Server{
		Handler: handler,
		Addr:    ":3000",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Server started on port 3000")

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Error", err)
	}
}

func parseJSON(r *http.Request) (Payload, error) {
	var payload Payload
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Payload{}, err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &payload); err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
