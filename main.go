package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/charmbracelet/log"
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

const PORT int = 3000

var logger = slog.New(log.New(os.Stderr))

// @title         Calculator API
// @version       1.0
// @description   This is a sample server for Calculator API
//
// @contact.name  Abdoulaye NDOYE
// @contact.url   https://github.com/NDOY3M4N
// @contact.email pa.ndoye@outlook.com
//
// @license.name  MIT
// @license.url   https://github.com/NDOY3M4N/calculator-api/blob/main/LICENSE
//
// @host          localhost:3000
// @BasePath      /api/v1
func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /add", addHandler)
	router.HandleFunc("POST /substract", substractHandler)
	router.HandleFunc("POST /multiply", multiplyHandler)
	router.HandleFunc("POST /divide", divideHandler)

	// Define a separate handler for the /scalar endpoint
	scalarHandler := http.StripPrefix(
		"/scalar",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL:       "./docs/swagger.json",
				CustomOptions: scalar.CustomOptions{PageTitle: "P4P1's API doc"},
				DarkMode:      true,
			})
			if err != nil {
				fmt.Printf("%v", err)
			}

			fmt.Fprintln(w, htmlContent)
		}),
	)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// Combine the v1 handler and the Swagger handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/scalar" || strings.HasPrefix(r.URL.Path, "/scalar") {
			scalarHandler.ServeHTTP(w, r)
		} else {
			v1.ServeHTTP(w, r)
		}
	})

	server := http.Server{
		Handler: Logging(logger, handler),
		Addr:    fmt.Sprintf(":%d", PORT),
	}

	logger.Info(fmt.Sprintf("Server started on port :%d", PORT))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Error", err)
	}
}
