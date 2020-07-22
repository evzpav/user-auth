package http

import (
	"net/http"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

func (h *handler) getLoginGoogle(w http.ResponseWriter, r *http.Request) {
	// if _, ok := h.alreadyLoggedIn(w, r); ok {
	// 	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	// 	return
	// }

	session, err := h.store.Get(r, googleSession)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	state := h.authService.GenerateToken()
	session.Values[googleCookie] = state
	if err := session.Save(r, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	googleSigninLink := h.authService.GetGoogleSigninLink(state)

	http.Redirect(w, r, googleSigninLink, http.StatusTemporaryRedirect)
}

func (h *handler) googleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := h.store.Get(r, googleSession)
	if err != nil {
		h.log.Error().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queryStateMap := r.URL.Query()
	queryState := queryStateMap.Get("state")

	sessionState, ok := session.Values[googleCookie]
	if !ok {
		h.log.Error().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sessionState.(string) != queryState {
		h.log.Info().Sendf("Invalid session state: retrieved: %s; Param: %s", sessionState, queryState)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	code := queryStateMap.Get("code")

	googleUser, err := h.authService.GetGoogleProfile(code)
	if err != nil {
		h.log.Error().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !googleUser.EmailVerified {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userService.FindByGoogleID(ctx, googleUser.Sub)
	if err != nil {
		h.log.Error().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		authUser := &domain.AuthUser{
			Email:    googleUser.Email,
			Password: googleUser.Sub,
			Name:     googleUser.Name,
			GoogleID: googleUser.Sub,
		}

		if err := h.authService.SignupWithGoogle(ctx, authUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	userProfile := profile{
		Profile: domain.Profile{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Address: user.Address,
			Phone:   user.Phone,
		},
		Errors: make(map[string]string),
	}

	h.writeTemplate(w, "profile", userProfile)

	// http.Redirect(w, r, "/profile", http.StatusSeeOther)

}
