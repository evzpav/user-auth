package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

const authSession string = "user_auth_session"
const authCookie string = "user_auth"
const googleSession string = "google_session"
const googleCookie string = "google_cookie"
const sessionLength int = 86400 * 7 // 1 week in seconds

var defaultSessionOptions = &sessions.Options{
	Path:     "/",
	HttpOnly: true,
	MaxAge:   sessionLength,
}

type profile struct {
	domain.Profile
	Errors  map[string]string
	Message string
}

type handler struct {
	userService     domain.UserService
	authService     domain.AuthService
	templateService domain.TemplateService
	store           *sessions.CookieStore
	log             log.Logger
}

func NewHandler(userService domain.UserService, authService domain.AuthService, templateService domain.TemplateService, sessionKey string, log log.Logger) http.Handler {
	handler := &handler{
		userService:     userService,
		authService:     authService,
		templateService: templateService,
		store:           sessions.NewCookieStore([]byte(sessionKey)),
		log:             log,
	}

	r := mux.NewRouter()
	r.Use(handler.logger())

	r.HandleFunc("/", redirectToLogin).Methods("GET")
	r.HandleFunc("/login", handler.getLogin).Methods("GET")
	r.HandleFunc("/login", handler.postLogin).Methods("POST")
	r.HandleFunc("/login/google", handler.getLoginGoogle).Methods("GET")
	r.HandleFunc("/login/google/auth", handler.googleAuth).Methods("GET")
	r.HandleFunc("/signup", handler.getSignup).Methods("GET")
	r.HandleFunc("/signup", handler.postSignup).Methods("POST")
	r.HandleFunc("/logout", handler.logout).Methods("GET")
	r.HandleFunc("/password/forgot", handler.getForgotPassword).Methods("GET")
	r.HandleFunc("/password/forgot", handler.postForgotPassword).Methods("POST")
	r.HandleFunc("/password/new", handler.getNewPassword).Methods("GET")
	r.HandleFunc("/password/new", handler.postNewPassword).Methods("POST")
	r.HandleFunc("/profile", handler.postProfile).Methods("POST")
	r.HandleFunc("/profile", handler.getProfile).Methods("GET")
	r.HandleFunc("/address", handler.getAddressSuggestion).Methods("GET")

	return r
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (h *handler) alreadyLoggedIn(w http.ResponseWriter, r *http.Request) (*domain.User, bool) {
	ctx := r.Context()
	session, err := h.store.Get(r, authSession)
	if err != nil {
		return nil, false
	}

	token, ok := session.Values[authCookie]
	if !ok {
		return nil, false
	}

	user, err := h.authService.AuthenticateToken(ctx, token.(string))
	if err != nil {
		return nil, false
	}

	session.Options = defaultSessionOptions
	session.Values[authCookie] = user.Token
	if err := session.Save(r, w); err != nil {
		return nil, false
	}

	return user, true
}

func (h *handler) writeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	loginTpl, err := h.templateService.RetrieveParsedTemplate(templateName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := loginTpl.Template.Execute(w, data); err != nil {
		h.log.Error().Sendf("Failed to execute template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) getSessionAndSetCookie(w http.ResponseWriter, r *http.Request, token, sessionName, cookieName string, options *sessions.Options) error {
	session, err := h.store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Options = options
	session.Values[cookieName] = token
	return session.Save(r, w)
}
