package http

import (
	"net/http"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

func (h *handler) getSignup(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.alreadyLoggedIn(w, r); ok {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.writeTemplate(w, "signup", nil)
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

	if err := h.getSessionAndSetCookie(w, r, authUser.Token); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.writeTemplate(w, "signup", authUser)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)

}
