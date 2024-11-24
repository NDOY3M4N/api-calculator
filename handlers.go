package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

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
		logger.Error("Division by zero",
			slog.Int("statusCode", http.StatusBadRequest),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Division by zero"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"result": payload.Number1 / payload.Number2})
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
