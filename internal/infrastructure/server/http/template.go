package http

import (
	"net/http"
	"time"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

type session struct {
	un           string
	lastActivity time.Time
}

var dbUsers = map[string]domain.User{}
var dbSessions = map[string]session{}
var dbSessionsCleaned time.Time

const sessionLength int = 30

func init() {
	// bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	// dbUsers["test@test.com"] = user{"test@test.com", bs, "James", "Bond"}
	dbSessionsCleaned = time.Now()
}

func (h *handler) writeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
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

func (h *handler) getLogin(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.writeTemplate(w, "login", nil)
}

func (h *handler) getSignup(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.writeTemplate(w, "signup", nil)
}

func (h *handler) postLogin(w http.ResponseWriter, r *http.Request) {
	authUser := domain.NewAuthUser(r.FormValue("email"), r.FormValue("password"))
	if !authUser.Validate() {
		h.writeTemplate(w, "login", authUser)
		return
	}

	err := h.authService.Authenticate(r.Context(), authUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.writeTemplate(w, "login", authUser)
		return
	}

	c := newCookie(authUser.Token, sessionLength)
	http.SetCookie(w, c)

	dbSessions[c.Value] = session{authUser.Email, time.Now()}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	return

}

func (h *handler) postSignup(w http.ResponseWriter, r *http.Request) {
	authUser := domain.NewAuthUser(r.FormValue("email"), r.FormValue("password"))
	if !authUser.Validate() {
		h.writeTemplate(w, "signup", authUser)
		return
	}

	if err := h.authService.Signup(r.Context(), authUser); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.writeTemplate(w, "signup", authUser)
		return
	}

	c := newCookie(authUser.Token, sessionLength)

	http.SetCookie(w, c)
	dbSessions[c.Value] = session{authUser.Email, time.Now()}

	http.Redirect(w, r, "/profile", http.StatusContinue)
	return

}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c, _ := r.Cookie("session")

	delete(dbSessions, c.Value)
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *handler) getProfile(w http.ResponseWriter, r *http.Request) {
	h.writeTemplate(w, "profile", nil)
}

func (h *handler) postProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("postprofile"))
}

// func (h *handler) putProfile(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("putprofile"))
// }

func (h *handler) getResetPassword(w http.ResponseWriter, r *http.Request) {
	h.writeTemplate(w, "reset_password", nil)
}

func (h *handler) postResetPassword(w http.ResponseWriter, r *http.Request) {
	h.writeTemplate(w, "email_sent", nil)
}

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
}
