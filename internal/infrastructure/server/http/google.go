package http

import (
	"net/http"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

func (h *handler) getLoginGoogle(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.alreadyLoggedIn(w, r); ok {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	state := h.authService.GenerateToken()

	err := h.getSessionAndSetCookie(w, r, state, googleSession, googleCookie, defaultSessionOptions)
	if err != nil {
		h.log.Info().Sendf("failed to get session and set cookie")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	googleSigninLink := h.authService.GetGoogleSigninLink(state)

	http.Redirect(w, r, googleSigninLink, http.StatusTemporaryRedirect)
}

func (h *handler) googleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if _, ok := h.alreadyLoggedIn(w, r); ok {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	session, err := h.store.Get(r, googleSession)
	if err != nil {
		h.log.Info().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	queryStateMap := r.URL.Query()
	queryState := queryStateMap.Get("state")

	sessionState, ok := session.Values[googleCookie]
	if !ok {
		h.log.Info().Send("failed to get session state")
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionStateStr, ok := sessionState.(string)
	if !ok {
		h.log.Info().Send("failed to cast session state")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sessionStateStr != queryState {
		h.log.Info().Sendf("Invalid session state: retrieved: %s; Param: %s", sessionState, queryState)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	code := queryStateMap.Get("code")

	googleUser, err := h.authService.GetGoogleProfile(code)
	if err != nil {
		h.log.Info().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !googleUser.EmailVerified || googleUser.Sub == "" {
		h.log.Info().Send("email not verified or googleID empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userService.FindByGoogleID(ctx, googleUser.Sub)
	if err != nil {
		h.log.Info().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authUser := &domain.AuthUser{
		Email:    googleUser.Email,
		Password: googleUser.Sub,
		Name:     googleUser.Name,
		GoogleID: googleUser.Sub,
	}

	if user == nil {
		user, err = h.authService.SignupWithGoogle(ctx, authUser)
		if err != nil {
			h.log.Info().Err(err).Send(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if user.Name == "" && googleUser.Name != "" {
		user.Name = googleUser.Name
	}

	user, err = h.authService.SetToken(ctx, user)
	if err != nil {
		h.log.Info().Err(err).Send(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.getSessionAndSetCookie(w, r, user.Token, authSession, authCookie, defaultSessionOptions); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.writeTemplate(w, "login", authUser)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)

}
