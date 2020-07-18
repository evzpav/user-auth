package http

import (
	"fmt"
	"net/http"
)

func (h *handler) getLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Login")
}
