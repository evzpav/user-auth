package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

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
	user, ok := h.alreadyLoggedIn(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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

func (h *handler) getAddressSuggestion(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	addressInput := queryParams.Get("q")

	if len(addressInput) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	suggestion, err := h.templateService.GetAddressSuggestion(addressInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(suggestion)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
