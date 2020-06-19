package core

import (
	"fmt"
	"net/http"
)

// healthcheckHandler is a handler that returns 200 to signify that the API is in good health for receiving requests.
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Print a debug message.
	fmt.Println("OK")

	w.WriteHeader(http.StatusOK)
}
