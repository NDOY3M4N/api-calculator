package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type Payload struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var req Payload
	if err := parseJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": req.Number1 + req.Number2})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", addHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":3000",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Server started on port 3000")

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Error", err)
	}
}

func parseJSON(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.Unmarshal(body, v)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
