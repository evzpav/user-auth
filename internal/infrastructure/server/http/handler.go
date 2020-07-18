package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

type handler struct {
	userService domain.UserService
	log         log.Logger
}

// NewHandler ...
func NewHandler(us domain.UserService, log log.Logger) http.Handler {
	handler := &handler{
		userService: us,
		log:         log,
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", handler.getLogin).Methods("GET")

	return r
}
