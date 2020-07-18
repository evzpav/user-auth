package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

type handler struct {
	userService     domain.UserService
	templateService domain.TemplateService
	log             log.Logger
}

// NewHandler ...
func NewHandler(us domain.UserService, ts domain.TemplateService, log log.Logger) http.Handler {
	handler := &handler{
		userService:     us,
		templateService: ts,
		log:             log,
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", handler.getLogin).Methods(http.MethodGet)
	r.HandleFunc("/login", handler.postLogin).Methods(http.MethodPost)
	r.HandleFunc("/signup", handler.getSignup).Methods(http.MethodGet)
	r.HandleFunc("/signup", handler.postSignup).Methods(http.MethodPost)
	r.HandleFunc("/profile", authorized(handler.getProfile)).Methods(http.MethodGet)
	r.HandleFunc("/profile", handler.postProfile).Methods(http.MethodPost)
	r.HandleFunc("/resetpassword", handler.resetPassword).Methods(http.MethodPost)
	r.HandleFunc("/logout", handler.logout)

	return r
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
	}
	_, ok = dbUsers[s.un]
	c.MaxAge = 30
	http.SetCookie(w, c)
	return ok
}

func authorized(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !alreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}
