package resp

import (
	"net/http"

	"github.com/go-chi/render"
)

type response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []Err       `json:"errors,omitempty"`
}

// Err represents a client-facing error.
type Err struct {
	Code          int         `json:"code"`
	Message       string      `json:"msg"`
	Data          interface{} `json:"data,omitempty"`
	InvalidFields interface{} `json:"invalidFields,omitempty"`
}

// Client-facing standard errors.
var (
	UnknownError      = Error(1, "unknown error")
	UnauthorizedError = Error(401, "unauthorized; check headers and token")
)

// OK returns an HTTP 200 response.
func OK(w http.ResponseWriter, r *http.Request, data ...interface{}) {
	resp := response{}

	// for i, v := range data {
	// Loop through each data object to check its type
	// data[i] = middleware.CleanDataForJSON(v)
	// }

	resp.Data = data
	if len(data) == 1 {
		resp.Data = data[0]
	}

	render.JSON(w, r, resp)
}

// NotFound returns an HTTP 404 response.
func NotFound(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusNotFound)
	render.JSON(w, r, resp)
}

// ServerError returns an HTTP 500 response.
func ServerError(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(w, r, resp)
}

// BadRequest returns an HTTP 400 response.
func BadRequest(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, resp)
}

// Unauthorized returns an HTTP 401 response.
func Unauthorized(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusUnauthorized)
	render.JSON(w, r, resp)
}

// Forbidden returns an HTTP 403 response.
func Forbidden(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusForbidden)
	render.JSON(w, r, resp)
}

// Error returns a client-facing error.
func Error(code int, message string) Err {
	return Err{
		Code:    code,
		Message: message,
	}
}

// ErrorData returns a client-facing error with data.
func ErrorData(code int, message string, data interface{}) Err {
	return Err{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// ErrorInvalidFields returns a client-facing error with invalid fields.
func ErrorInvalidFields(code int, message string, invalidFields interface{}) Err {
	return Err{
		Code:          code,
		Message:       message,
		InvalidFields: invalidFields,
	}
}

// ErrorCheckData returns a client-facing error 0.
func ErrorCheckData(data interface{}) Err {
	return Err{
		Code:    0,
		Message: "check data",
		Data:    data,
	}
}
