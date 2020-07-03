package resp

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/joinimpact/api/pkg/scopes"
	"github.com/liip/sheriff"
)

type response struct {
	Data   interface{} `json:"data,omitempty" groups:"user"`
	Errors []Err       `json:"errors,omitempty"`
}

// Err represents a client-facing error.
type Err struct {
	Code          int         `json:"code"`
	Message       string      `json:"msg"`
	Data          interface{} `json:"data,omitempty"`
	InvalidFields []string    `json:"invalidFields,omitempty"`
}

// Client-facing standard errors.
var (
	UnknownError      = Error(500, "unknown error")
	UnauthorizedError = Error(401, "unauthorized; check headers and token")
)

var defaultOptions = sheriff.Options{
	Groups: []string{"user"},
}

// marshal marshals the data and writes it to the response writer.
func marshal(w http.ResponseWriter, r *http.Request, groups []string, data interface{}) {
	options := defaultOptions
	if groups != nil {
		options.Groups = append(options.Groups, groups...)
	}

	// Marshal the data with scoping.
	data = scopes.MarshalFromContext(r.Context(), data)

	render.JSON(w, r, data)
}

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

	marshal(w, r, nil, resp)
}

// NotFound returns an HTTP 404 response.
func NotFound(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusNotFound)
	marshal(w, r, nil, resp)
}

// ServerError returns an HTTP 500 response.
func ServerError(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusInternalServerError)
	marshal(w, r, nil, resp)
}

// BadRequest returns an HTTP 400 response.
func BadRequest(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusBadRequest)
	marshal(w, r, nil, resp)
}

// Unauthorized returns an HTTP 401 response.
func Unauthorized(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusUnauthorized)
	marshal(w, r, nil, resp)
}

// Forbidden returns an HTTP 403 response.
func Forbidden(w http.ResponseWriter, r *http.Request, errors ...Err) {
	resp := response{}
	resp.Errors = errors

	w.WriteHeader(http.StatusForbidden)
	marshal(w, r, nil, resp)
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
func ErrorInvalidFields(code int, message string, invalidFields []string) Err {
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
