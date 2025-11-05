package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/diagnosis/luxsuv-api-v2/internal/helper"
	"github.com/diagnosis/luxsuv-api-v2/internal/secure"
	"github.com/diagnosis/luxsuv-api-v2/internal/store"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserStore    store.UserStore
	Signer       secure.Signer
	RefreshStore store.RefreshStore
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	//--1) parse & validate
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) //1mb
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	email := strings.ToLower(strings.TrimSpace(body.Email))
	pw := strings.TrimSpace(body.Password)
	if len(email) < 4 || len(pw) < 8 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	//2) find user
	u, err := h.UserStore.GetByEmail(ctx, email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !u.IsActive {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	//3) validate password
	if !secure.VerifyPassword(pw, u.PasswordHash) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	//4) mint access JWT
	access_token, _, err := h.Signer.MintAccess(u.ID, helper.DeferOrString(u.Role, "rider"))
	if err != nil {
		http.Error(w, "token mint failed", http.StatusInternalServerError)
		return
	}

	//5) create refresh 7 days
	ua := r.UserAgent()
	ip := helper.ClientIP(r)
	refreshToken := 

}
