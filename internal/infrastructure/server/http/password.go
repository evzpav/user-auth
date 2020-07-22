package http

import (
	"net/http"
	"strings"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

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

func (h *handler) getNewPassword(w http.ResponseWriter, r *http.Request) {
	var reply profile
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
	var reply profile
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
