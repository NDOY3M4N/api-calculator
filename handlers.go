package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"

	"github.com/NDOY3M4N/api-calculator/repository"
)

var (
	ErrMissingBody  = errors.New("missing request body")
	ErrDividyByZero = errors.New("division by zero is prohibited")
	ErrLengthSum    = errors.New("provide at least 2 numbers")
)

type Payload struct {
	Number1 float64 `json:"number1" example:"6"`
	Number2 float64 `json:"number2" example:"9"`
}

type PayloadSum []float64

type APIError struct {
	Error string `json:"error"`
}

type APISuccess struct {
	Result float64 `json:"result"`
}

type PayloadLogin struct {
	Pseudo string `json:"pseudo" example:"p4p1"`
}

type APILoginSuccess struct {
	Token string `json:"token"`
}

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) http.HandlerFunc {
	isAuth := IsAuthenticated(h.repo)

	router.HandleFunc("POST /login", h.loginHandler)

	router.HandleFunc("POST /add", isAuth(h.addHandler))
	router.HandleFunc("POST /sum", isAuth(h.sumHandler))
	router.HandleFunc("POST /substract", isAuth(h.substractHandler))
	router.HandleFunc("POST /multiply", isAuth(h.multiplyHandler))
	router.HandleFunc("POST /divide", isAuth(h.divideHandler))

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

// Login
//
// @summary Login
// @description Log the user
// @tags User
// @accept json
// @produce json
// @param payload body PayloadLogin true "Field needed for login"
// @success 200 {object} APILoginSuccess
// @failure 400 {object} APIError
// @router /login [post]
func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var payload PayloadLogin
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	if payload.Pseudo == "" {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("pseudo should not be empty"))
		return
	}

	user, err := h.repo.FindUserByPseudo(payload.Pseudo)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("error finding user by pseudo"))
		return
	}

	token, err := GenerateToken(int(user.Id))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("error generating token"))
		return
	}

	reqID := r.Context().Value(requestIDKey).(string)

	logger.Info("Request successful",
		slog.Int("statusCode", http.StatusOK),
		slog.String("remoteAddr", r.RemoteAddr),
		slog.Group("request",
			slog.String("id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		),
	)

	encodeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// Add two numbers
//
// @summary Add two numbers
// @description Add two numbers together
// @tags Math
// @accept json
// @produce json
// @param payload body Payload true "Numbers needed for the operation"
// @Security BearerAuth
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /add [post]
func (h *Handler) addHandler(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	userID := r.Context().Value(userIDKey).(int)

	result := payload.Number1 + payload.Number2
	param := repository.AddOperationParams{
		Inputs: []float64{float64(payload.Number1), float64(payload.Number2)},
		Type:   repository.TypeAdd,
		Result: float64(result),
		UserId: userID,
	}

	if err := h.repo.AddOperation(param); err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}

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
// @Security BearerAuth
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /sum [post]
func (h *Handler) sumHandler(w http.ResponseWriter, r *http.Request) {
	var payload PayloadSum
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	if len(payload) < 2 {
		writeError(w, r, http.StatusBadRequest, ErrLengthSum)
		return
	}

	var result float64
	for _, num := range payload {
		result += num
	}

	userID := r.Context().Value(userIDKey).(int)

	param := repository.AddOperationParams{
		Inputs: payload,
		Type:   repository.TypeSum,
		Result: float64(result),
		UserId: userID,
	}

	if err := h.repo.AddOperation(param); err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
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
// @Security BearerAuth
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /substract [post]
func (h *Handler) substractHandler(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	result := payload.Number1 - payload.Number2

	userID := r.Context().Value(userIDKey).(int)

	param := repository.AddOperationParams{
		Inputs: []float64{float64(payload.Number1), float64(payload.Number2)},
		Type:   repository.TypeSubstract,
		Result: float64(result),
		UserId: userID,
	}

	if err := h.repo.AddOperation(param); err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}

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
// @Security BearerAuth
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /multiply [post]
func (h *Handler) multiplyHandler(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	result := payload.Number1 * payload.Number2

	userID := r.Context().Value(userIDKey).(int)

	param := repository.AddOperationParams{
		Inputs: []float64{float64(payload.Number1), float64(payload.Number2)},
		Type:   repository.TypeMultiply,
		Result: float64(result),
		UserId: userID,
	}

	if err := h.repo.AddOperation(param); err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}

	writeSuccess(w, r, http.StatusOK, result)
}

// divideHandler Foo
//
// @summary Divide two numbers
// @description Divide two numbers together
// @tags Math
// @accept json
// @produce json
// @param payload body Payload true "Numbers needed for the operation"
// @Security BearerAuth
// @success 200 {object} APISuccess
// @failure 400 {object} APIError
// @router /divide [post]
func (h *Handler) divideHandler(w http.ResponseWriter, r *http.Request) {
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

	userID := r.Context().Value(userIDKey).(int)

	param := repository.AddOperationParams{
		Inputs: []float64{float64(payload.Number1), float64(payload.Number2)},
		Type:   repository.TypeDivide,
		Result: float64(result),
		UserId: userID,
	}

	if err := h.repo.AddOperation(param); err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}

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

func writeSuccess(w http.ResponseWriter, r *http.Request, statusCode int, payload float64) error {
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
