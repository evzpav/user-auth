package http

import (
	"net/http"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

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
	if _, ok := h.alreadyLoggedIn(w, r); ok {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.writeTemplate(w, "login", nil)
}

func (h *handler) getSignup(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.alreadyLoggedIn(w, r); ok {
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

	// err := h.authService.Authenticate(r.Context(), authUser)
	user, err := h.authService.Authenticate2(r.Context(), authUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.writeTemplate(w, "login", authUser)
		return
	}

	c := newCookie(user.Token, sessionLength)
	http.SetCookie(w, c)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *handler) postSignup(w http.ResponseWriter, r *http.Request) {
	authUser := domain.NewAuthUser(r.FormValue("email"), r.FormValue("password"))
	if !authUser.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		h.writeTemplate(w, "signup", authUser)
		return
	}

	if err := h.authService.Signup(r.Context(), authUser); err != nil {
		w.WriteHeader(http.StatusForbidden)
		h.writeTemplate(w, "signup", authUser)
		return
	}

	c := newCookie(authUser.Token, sessionLength)
	http.SetCookie(w, c)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)

}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	// if !h.alreadyLoggedIn(w, r) {
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }

	if _, ok := h.alreadyLoggedIn(w, r); !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	c, err := r.Cookie(cookieName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	delete(h.session, c.Value)
	//TODO delete token form user

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *handler) getProfile(w http.ResponseWriter, r *http.Request) {
	// user, ok := h.alreadyLoggedIn(w, r)
	// if !ok {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }

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
	authUser := domain.NewAuthUser(r.FormValue("email"), "")

	if !authUser.ValidateEmail() {
		w.WriteHeader(http.StatusBadRequest)
		h.writeTemplate(w, "reset_password", authUser)
		return
	}

	message := "message"

	if err := h.authService.SendEmail(r.Context(), message, authUser.Email); err != nil {
		errorMsg := "failed to send email"
		h.log.Error().Err(err).Sendf(errorMsg)
		authUser.Errors["Credentials"] = errorMsg

		w.WriteHeader(http.StatusBadRequest)
		h.writeTemplate(w, "reset_password", authUser)
		return
	}

	h.writeTemplate(w, "email_sent", nil)
}
