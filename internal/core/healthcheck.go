package core

import "net/http"

// healthcheckHandler is a handler that returns 200 to signify that the API is in good health for receiving requests.
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
