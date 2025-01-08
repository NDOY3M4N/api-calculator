package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/charmbracelet/log"

	"github.com/NDOY3M4N/api-calculator/ratelimit"
)

const (
	port       int   = 3000
	bucketSize int64 = 5
	bucketRate int64 = 3
)

var logger = slog.New(log.New(os.Stderr))

// @title         Calculator API
// @version       1.0
// @description   This is a simple server for Calculator API
//
// @contact.name  Abdoulaye NDOYE
// @contact.url   https://github.com/NDOY3M4N
// @contact.email pa.ndoye@outlook.com
//
// @license.name  MIT
// @license.url   https://github.com/NDOY3M4N/api-calculator/blob/main/LICENSE.md
//
// @servers.url http://localhost:3000/api/v1
// @servers.description Development server
func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /add", addHandler)
	router.HandleFunc("POST /sum", sumHandler)
	router.HandleFunc("POST /substract", substractHandler)
	router.HandleFunc("POST /multiply", multiplyHandler)
	router.HandleFunc("POST /divide", divideHandler)

	// Define a separate handler for the /scalar endpoint
	scalarHandler := http.StripPrefix(
		"/docs",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL: "./docs/swagger.json",
				CustomOptions: scalar.CustomOptions{
					PageTitle: "P4P1's Calculator API doc",
				},
				// Layout:   scalar.LayoutClassic,
				DarkMode: true,
			})
			if err != nil {
				fmt.Printf("%v", err)
			}

			fmt.Fprintln(w, htmlContent)
		}),
	)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// Combine the v1 handler and the Scalar handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/docs" || strings.HasPrefix(r.URL.Path, "/docs") {
			scalarHandler.ServeHTTP(w, r)
		} else {
			v1.ServeHTTP(w, r)
		}
	})

	bucket := ratelimit.NewTokenBucket(bucketSize, bucketRate)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bucket.Start(ctx)

	stack := CreateStack(AddRequestId, Logger, RateLimit(bucket))
	server := http.Server{
		Handler: stack(handler),
		Addr:    fmt.Sprintf(":%d", port),
	}

	logger.Info(fmt.Sprintf("Server started on port :%d", port))
	logger.Info(fmt.Sprintf("API documentation available on http://localhost:%d/docs", port))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Error", err.Error(), nil)
	}
}
