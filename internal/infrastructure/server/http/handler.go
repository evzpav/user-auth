package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

const cookieName string = "user_auth_session"
const sessionLength int = 60

type handler struct {
	userService     domain.UserService
	authService     domain.AuthService
	templateService domain.TemplateService
	session         map[string]*domain.User
	log             log.Logger
}

// NewHandler ...
func NewHandler(userService domain.UserService, authService domain.AuthService, templateService domain.TemplateService, log log.Logger) http.Handler {
	handler := &handler{
		userService:     userService,
		authService:     authService,
		templateService: templateService,
		session:         make(map[string]*domain.User),
		log:             log,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", redirectToLogin).Methods("GET")
	r.HandleFunc("/login", handler.getLogin).Methods("GET")
	r.HandleFunc("/login", handler.postLogin).Methods("POST")
	r.HandleFunc("/signup", handler.getSignup).Methods("GET")
	r.HandleFunc("/signup", handler.postSignup).Methods("POST")
	r.HandleFunc("/profile", handler.getProfile).Methods("GET")
	r.HandleFunc("/profile", handler.postProfile).Methods("POST")
	r.HandleFunc("/resetpassword", handler.getResetPassword).Methods("GET")
	r.HandleFunc("/resetpassword", handler.postResetPassword).Methods("POST")
	r.HandleFunc("/logout", handler.logout)

	return r
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// func (h *handler) alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
// 	c, err := r.Cookie(cookieName)
// 	if err != nil {
// 		return false
// 	}

// 	if c.MaxAge <= 0 {
// 		return false
// 	}

// 	err = h.authService.AuthenticateToken(r.Context(), c.Value)
// 	if err != nil {
// 		return false
// 	}

// 	c.MaxAge = sessionLength
// 	http.SetCookie(w, c)

// 	return true
// }

func (h *handler) alreadyLoggedIn(w http.ResponseWriter, r *http.Request) (*domain.User, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil, false
	}

	if c.MaxAge <= 0 {
		return nil, false
	}

	token := c.Value

	user, ok := h.session[token]
	if ok {
		return user, true
	}

	user, err = h.authService.AuthenticateToken(r.Context(), c.Value)
	if err != nil {
		return nil, false
	}

	h.session[user.Token] = user

	newCookie := newCookie(user.Token, sessionLength)
	http.SetCookie(w, newCookie)

	return user, true
}

func (h *handler) authorized(hl http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if !h.alreadyLoggedIn(w, r) {
		// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
		// 	return
		// }
		if _, ok := h.alreadyLoggedIn(w, r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		hl.ServeHTTP(w, r)
	})
}

func newCookie(token string, age int) *http.Cookie {
	c := &http.Cookie{
		Name:  cookieName,
		Value: token,
	}
	c.MaxAge = age
	return c
}
