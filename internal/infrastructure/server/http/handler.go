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
	authService     domain.AuthService
	templateService domain.TemplateService
	log             log.Logger
}

// NewHandler ...
func NewHandler(userService domain.UserService, authService domain.AuthService, templateService domain.TemplateService, log log.Logger) http.Handler {
	handler := &handler{
		userService:     userService,
		authService:     authService,
		templateService: templateService,
		log:             log,
	}

	r := mux.NewRouter()
	
	r.HandleFunc("/", redirectToLogin).Methods("GET")
	r.HandleFunc("/login", handler.getLogin).Methods("GET")
	r.HandleFunc("/login", handler.postLogin).Methods("POST")
	r.HandleFunc("/signup", handler.getSignup).Methods("GET")
	r.HandleFunc("/signup", handler.postSignup).Methods("POST")
	r.HandleFunc("/profile", authorized(handler.getProfile)).Methods("GET")
	r.HandleFunc("/profile", handler.postProfile).Methods("POST")
	r.HandleFunc("/resetpassword", handler.getResetPassword).Methods("GET")
	r.HandleFunc("/resetpassword", handler.postResetPassword).Methods("POST")
	r.HandleFunc("/logout", handler.logout)

	return r
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
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

func newCookie(token string, age int) *http.Cookie {
	c := &http.Cookie{
		Name:  "session",
		Value: token,
	}
	c.MaxAge = age
	return c
}

