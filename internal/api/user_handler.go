package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/jingxinwangdev/go-prject/internal/store"
	"github.com/jingxinwangdev/go-prject/internal/utils"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) ValidateUserRequest(req RegisterUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	// https://pkg.go.dev/regexp
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email address")
	}
	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("Error decoding user request: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}
	err = uh.ValidateUserRequest(req)
	if err != nil {
		uh.logger.Printf("Error validating user request: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// Deal with password
	/**
	* 整个流程是：
	* &store.User{...} → PasswordHash 自动初始化为零值（空 struct）
	* .Set(req.Password) → 把 plaintext 和 hash 填上真实值
	* CreateUser(user) → 把 hash 存入数据库
	**/
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("Error setting password: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to set password"})
		return
	}
	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("Error creating user: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create user"})
		return
	}
	utils.WriteJsonResponse(w, http.StatusCreated, utils.Envelope{"user": user})
}
