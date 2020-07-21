package http

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"gitlab.com/evzpav/user-auth/internal/domain"
)

type profile struct {
	domain.Profile
	Errors map[string]string
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

	user, err := h.authService.Authenticate(r.Context(), authUser)
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
	user, ok := h.alreadyLoggedIn(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	delete(h.sessions, user.Token)

	c := newCookie("", -1)
	http.SetCookie(w, c)

	user.Token = ""

	if err := h.userService.Update(r.Context(), user); err != nil {
		h.log.Error().Err(err).Sendf("failed to update user token")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (h *handler) getProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := h.alreadyLoggedIn(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	prof := profile{
		Profile: domain.Profile{
			ID:      user.ID,
			Email:   user.Email,
			Address: user.Address,
			Phone:   user.Phone,
			Name:    user.Name,
		},
		Errors: make(map[string]string),
	}

	h.writeTemplate(w, "profile", prof)
}

func (h *handler) postProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userProfile := profile{
		Profile: domain.Profile{
			ID:      id,
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Address: r.FormValue("address"),
			Phone:   r.FormValue("phone"),
		},
		Errors: make(map[string]string),
	}

	if err := userProfile.Validate(); err != nil {
		h.log.Error().Err(err).Sendf("invalid user attributes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userService.FindByID(ctx, id)
	if err != nil {
		h.log.Error().Err(err).Sendf("failed to find user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Name = userProfile.Name
	user.Email = userProfile.Email
	user.Address = userProfile.Address
	user.Phone = userProfile.Phone

	if err := h.userService.Update(ctx, user); err != nil {
		h.log.Error().Err(err).Sendf("failed to update user profile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.writeTemplate(w, "profile", userProfile)
}

func (h *handler) getForgotPassword(w http.ResponseWriter, r *http.Request) {
	h.writeTemplate(w, "forgot_password", nil)
}

func (h *handler) postForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := domain.NewAuthUser(r.FormValue("email"), "")

	if !authUser.ValidateEmail() {
		w.WriteHeader(http.StatusBadRequest)
		h.writeTemplate(w, "forgot_password", authUser)
		return
	}

	token, err := h.authService.SetUserRecoveryToken(ctx, authUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		h.writeTemplate(w, "email_sent", nil) // to avoid bruteforce
		return
	}

	authUser.RecoveryToken = token

	go h.authService.SendResetPasswordLink(ctx, authUser)

	h.writeTemplate(w, "email_sent", nil)
}

type replyMessage struct {
	Errors  map[string]string
	Message string
}

func (h *handler) getNewPassword(w http.ResponseWriter, r *http.Request) {
	var reply replyMessage
	reply.Errors = make(map[string]string)

	urlValues := r.URL.Query()
	_, ok := urlValues["token"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		reply.Errors["Link"] = "invalid link"
		h.writeTemplate(w, "new_password", reply)
		return
	}

	h.writeTemplate(w, "new_password", reply)

}

func (h *handler) postNewPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reply replyMessage
	reply.Errors = make(map[string]string)

	authUser := domain.NewAuthUser("", r.FormValue("password"))
	if !authUser.ValidatePassword() {
		reply.Errors["Link"] = "invalid password"
		h.writeTemplate(w, "new_password", reply)
		return
	}

	token := r.FormValue("token")
	if strings.TrimSpace(token) == "" {
		w.WriteHeader(http.StatusBadRequest)
		reply.Errors["Link"] = "invalid link"
		h.writeTemplate(w, "new_password", reply)
		return
	}

	user, err := h.userService.FindByRecoveryToken(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		reply.Errors["Link"] = "invalid link"
		h.writeTemplate(w, "new_password", reply)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		reply.Errors["Link"] = "invalid link"
		h.writeTemplate(w, "new_password", reply)
		return
	}

	if err := h.authService.SetNewPassword(ctx, user, authUser.Password); err != nil {
		h.log.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		reply.Errors["Link"] = "failed to change password"
		h.writeTemplate(w, "new_password", reply)
		return
	}

	reply.Errors = nil
	reply.Message = "password changed"
	h.writeTemplate(w, "new_password", reply)
}

func (h *handler) loginWithThirdParty(w http.ResponseWriter, r *http.Request) {
	key := "google_cookie"
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	h.authService.GoogleAuthentication(w, r, store)
}

func (h *handler) loginWithThirdPartyCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return
	}
	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(w, user)
}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
