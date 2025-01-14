package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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
// @securitydefinitions.bearerauth BearerAuth
// @name Authorization
// @in header
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
	handler := NewHandler(repo).RegisterRoutes(router)

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

func initStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		logger.Error("DB initialization", slog.String("message", err.Error()))
		os.Exit(1)
	}

	logger.Info("DB successfully connected")
}
