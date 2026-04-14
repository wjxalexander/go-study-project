package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jingxinwangdev/go-prject/internal/store"
	"github.com/jingxinwangdev/go-prject/internal/tokens"
	"github.com/jingxinwangdev/go-prject/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{tokenStore: tokenStore, userStore: userStore, logger: logger}
}

func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var request createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		th.logger.Printf("Error decoding request: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}
	// go the user
	user, err := th.userStore.GetUserByUsername(request.Username)
	if err != nil || user == nil {
		th.logger.Printf("Error getting user: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	// password match
	passwordMatch, err := user.PasswordHash.Compare(request.Password)
	if err != nil {
		th.logger.Printf("Error comparing password: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if !passwordMatch {
		th.logger.Printf("Invalid password")
		utils.WriteJsonResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid password or username"})
		return
	}
	token, err := th.tokenStore.CreateToken(user.ID, time.Hour*24*30, tokens.ScopeAuthentication)
	if err != nil {
		th.logger.Printf("Error creating token: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJsonResponse(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
