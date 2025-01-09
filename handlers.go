package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrMissingBody  = errors.New("missing request body")
	ErrDividyByZero = errors.New("division by zero is prohibited")
	ErrLengthSum    = errors.New("provide at least 2 numbers")
)

type number float64

type Payload struct {
	Number1 number `json:"number1" example:"6"`
	Number2 number `json:"number2" example:"9"`
}

type PayloadSum []number

type APIError struct {
	Error string `json:"error"`
}

type APISuccess struct {
	Result number `json:"result"`
}

// Add two numbers
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
	var payload Payload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	result := payload.Number1 + payload.Number2
	writeSuccess(w, r, http.StatusOK, result)
}

// Sum numbers
//
// @summary Sum numbers
// @description Add all numbers in an array
// @tags Math
// @accept json
// @produce json
// @param payload body PayloadSum true "Array of numbers needed for the operation"
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /sum [post]
func sumHandler(w http.ResponseWriter, r *http.Request) {
	var payload PayloadSum
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	if len(payload) < 2 {
		writeError(w, r, http.StatusBadRequest, ErrLengthSum)
		return
	}

	var result number
	for _, num := range payload {
		result += num
	}

	writeSuccess(w, r, http.StatusOK, result)
}

// Substract two numbers
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
	var payload Payload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	result := payload.Number1 - payload.Number2
	writeSuccess(w, r, http.StatusOK, result)
}

// Multiply two numbers
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
	var payload Payload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	result := payload.Number1 * payload.Number2
	encodeJSON(w, http.StatusOK, result)
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
	var payload Payload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	if payload.Number2 == 0 {
		writeError(w, r, http.StatusBadRequest, ErrDividyByZero)
		return
	}

	result := payload.Number1 / payload.Number2
	writeSuccess(w, r, http.StatusOK, result)
}

func decodeJSON(r *http.Request, payload any) error {
	if r.ContentLength == 0 {
		return ErrMissingBody
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func encodeJSON(w http.ResponseWriter, statusCode int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}

func writeSuccess(w http.ResponseWriter, r *http.Request, statusCode int, payload number) error {
	reqID := r.Context().Value(requestIDKey).(string)

	logger.Info("Request successful",
		slog.Int("statusCode", statusCode),
		slog.String("remoteAddr", r.RemoteAddr),
		slog.Group("request",
			slog.String("id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		),
	)

	return encodeJSON(w, statusCode, APISuccess{payload})
}

func writeError(w http.ResponseWriter, r *http.Request, statusCode int, err error) error {
	reqID := r.Context().Value(requestIDKey).(string)

	logger.Error(err.Error(),
		slog.Int("statusCode", statusCode),
		slog.String("remoteAddr", r.RemoteAddr),
		slog.Group("request",
			slog.String("id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		),
	)

	return encodeJSON(w, statusCode, APIError{err.Error()})
}
