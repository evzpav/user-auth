package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

const authCookie string = "user_auth_session"
const sessionLength int = 3600 // 1 hour in seconds

type session struct {
	*domain.User
	ExpireTime time.Time
}

func newSession(user *domain.User) *session {
	return &session{
		User:       user,
		ExpireTime: time.Now().UTC().Add(time.Second * time.Duration(sessionLength)),
	}
}

type handler struct {
	userService     domain.UserService
	authService     domain.AuthService
	templateService domain.TemplateService
	sessions        map[string]*session
	log             log.Logger
}

// NewHandler ...
func NewHandler(userService domain.UserService, authService domain.AuthService, templateService domain.TemplateService, log log.Logger) http.Handler {
	handler := &handler{
		userService:     userService,
		authService:     authService,
		templateService: templateService,
		sessions:        make(map[string]*session),
		log:             log,
	}

	r := mux.NewRouter()
	r.Use(handler.logger())

	r.HandleFunc("/", redirectToLogin).Methods("GET")
	r.HandleFunc("/login", handler.getLogin).Methods("GET")
	r.HandleFunc("/login", handler.postLogin).Methods("POST")
	r.HandleFunc("/login/{provider}", handler.loginWithThirdParty).Methods("GET")
	r.HandleFunc("/login/{provider}/callback", handler.loginWithThirdPartyCallback).Methods("GET")
	r.HandleFunc("/signup", handler.getSignup).Methods("GET")
	r.HandleFunc("/signup", handler.postSignup).Methods("POST")
	r.HandleFunc("/logout", handler.logout).Methods("GET")
	r.HandleFunc("/password/forgot", handler.getForgotPassword).Methods("GET")
	r.HandleFunc("/password/forgot", handler.postForgotPassword).Methods("POST")
	r.HandleFunc("/password/new", handler.getNewPassword).Methods("GET")
	r.HandleFunc("/password/new", handler.postNewPassword).Methods("POST")
	r.HandleFunc("/profile", handler.postProfile).Methods("POST")
	r.HandleFunc("/profile", handler.getProfile).Methods("GET")

	return r
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (h *handler) alreadyLoggedIn(w http.ResponseWriter, r *http.Request) (*domain.User, bool) {
	h.cleanExpiredSessions()

	c, err := r.Cookie(authCookie)
	if err != nil {
		return nil, false
	}

	token := c.Value

	// session, ok := h.sessions[token]
	// if ok {
	// 	c.MaxAge = sessionLength
	// 	http.SetCookie(w, c)
	// 	return session.User, true
	// }

	user, err := h.authService.AuthenticateToken(r.Context(), token)
	if err != nil {
		return nil, false
	}

	// h.sessions[user.Token] = newSession(user)

	newCookie := newCookie(user.Token, sessionLength)
	http.SetCookie(w, newCookie)

	return user, true
}

func (h *handler) authorized(hl http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := h.alreadyLoggedIn(w, r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		hl.ServeHTTP(w, r)
	})
}

func newCookie(token string, age int) *http.Cookie {
	return &http.Cookie{
		Name:   authCookie,
		Value:  token,
		MaxAge: age,
	}
}

func (h *handler) cleanExpiredSessions() {
	for k, s := range h.sessions {
		if time.Now().UTC().After(s.ExpireTime) {
			delete(h.sessions, k)
		}
	}
}
