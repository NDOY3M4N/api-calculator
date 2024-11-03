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
	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 + payload.Number2})
}

func substractHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 - payload.Number2})
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := parseJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 * payload.Number2})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", addHandler)
	mux.HandleFunc("POST /substract", substractHandler)
	mux.HandleFunc("POST /multiply", multiplyHandler)

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
