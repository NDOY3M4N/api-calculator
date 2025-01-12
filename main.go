package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/charmbracelet/log"

	"github.com/NDOY3M4N/api-calculator/ratelimit"
	"github.com/NDOY3M4N/api-calculator/repository"
)

const (
	port       int   = 3000
	bucketSize int64 = 5
	bucketRate int64 = 2
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
	db, err := NewDatabase(envs.DBString)
	if err != nil {
		panic(err)
	}

	initStorage(db)
	defer db.Close()

	repo := repository.New(db)

	router := http.NewServeMux()
	handler := registerRoutes(router, repo)

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
		logger.Error("HTTP server initialization", slog.String("message", err.Error()))
		os.Exit(1)
	}
}

func registerRoutes(router *http.ServeMux, repo *repository.Repository) http.HandlerFunc {
	isAuth := IsAuthenticated(repo)

	router.HandleFunc("POST /login", loginHandler(repo))
	router.HandleFunc("POST /add", isAuth(addHandler))
	router.HandleFunc("POST /sum", isAuth(sumHandler))
	router.HandleFunc("POST /substract", isAuth(substractHandler))
	router.HandleFunc("POST /multiply", isAuth(multiplyHandler))
	router.HandleFunc("POST /divide", isAuth(divideHandler))

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

	return handler
}

func initStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		logger.Error("DB initialization", slog.String("message", err.Error()))
		os.Exit(1)
	}

	logger.Info("DB successfully connected")
}
