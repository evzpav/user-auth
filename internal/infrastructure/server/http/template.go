package http

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Address struct {
	Street1 string `form:"label=Street;placeholder=123 Sample St"`
	Street2 string `form:"label=Street (cont);placeholder=Apt 123"`
	City    string
	State   string `form:"footer=Or your Province"`
	Zip     string `form:"label=Postal Code"`
	Country string
}

type user struct {
	UserName string
	Password []byte
	First    string
	Last     string
}

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

func (h *handler) writeTemplate(w http.ResponseWriter, templateName string) {
	w.Header().Set("Content-Type", "text/html")
	loginTpl, err := h.templateService.RetrieveParsedTemplate(templateName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := loginTpl.Template.Execute(w, nil); err != nil {
		h.log.Error().Sendf("Failed to execute template: %v", err)
		return
	}
}

func (h *handler) getLogin(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.writeTemplate(w, "login")
}

func (h *handler) getSignup(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.writeTemplate(w, "signup")
}

func (h *handler) postLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	u, ok := dbUsers[email]
	if !ok {
		http.Error(w, "Email and/or password do not match", http.StatusForbidden)
		return
	}

	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		http.Error(w, "Email and/or password do not match", http.StatusForbidden)
		return
	}

	sID := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	dbSessions[c.Value] = session{email, time.Now()}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	return

}

func (h *handler) postSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if _, ok := dbUsers[email]; ok {
		http.Error(w, "Username already taken", http.StatusForbidden)
		return
	}

	sID := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	dbSessions[c.Value] = session{email, time.Now()}

	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	u := &domain.User{
		Email:    email,
		Password: bs,
	}

	dbUsers[email] = *u

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
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
	h.writeTemplate(w, "profile")
}

func (h *handler) postProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("postprofile"))
}

// func (h *handler) putProfile(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("putprofile"))
// }

func (h *handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("resetpassword"))
}

// func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
// 	c, err := req.Cookie("session")
// 	if err != nil {
// 		return false
// 	}
// 	s, ok := dbSessions[c.Value]
// 	if ok {
// 		s.lastActivity = time.Now()
// 		dbSessions[c.Value] = s
// 	}
// 	_, ok = dbUsers[s.un]
// 	c.MaxAge = 30
// 	http.SetCookie(w, c)
// 	return ok
// }

// func authorized(h http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !alreadyLoggedIn(w, r) {
// 			http.Redirect(w, r, "/", http.StatusSeeOther)
// 			return
// 		}
// 		h.ServeHTTP(w, r)
// 	})
// }

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
}
