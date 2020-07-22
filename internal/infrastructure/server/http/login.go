package http

import (
	"net/http"

	"github.com/gorilla/sessions"
	"gitlab.com/evzpav/user-auth/internal/domain"
)

func (h *handler) getLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.alreadyLoggedIn(w, r); ok {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.writeTemplate(w, "login", nil)
}

func (h *handler) postLogin(w http.ResponseWriter, r *http.Request) {
	authUser := domain.NewAuthUser(r.FormValue("email"), r.FormValue("password"))
	if !authUser.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		h.writeTemplate(w, "login", authUser)
		return
	}

	user, err := h.authService.Authenticate(r.Context(), authUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.writeTemplate(w, "login", authUser)
		return
	}

	if err := h.getSessionAndSetCookie(w, r, user.Token); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.writeTemplate(w, "login", authUser)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	user, ok := h.alreadyLoggedIn(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, err := h.store.Get(r, authSession)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	if err := session.Save(r, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Token = ""

	if err := h.userService.Update(r.Context(), user); err != nil {
		h.log.Error().Err(err).Sendf("failed to update user token")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
