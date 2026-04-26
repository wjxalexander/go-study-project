package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jingxinwangdev/go-prject/internal/store"
	"github.com/jingxinwangdev/go-prject/internal/tokens"
	"github.com/jingxinwangdev/go-prject/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

// https://pkg.go.dev/context
// https://gobyexample.com/context
func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("user not found in context")
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// within this anonymous function, we can access the request and response writer
		// we can inject the user into the request context
		w.Header().Set("Vary", "Authorization")
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorization, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJsonResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid authorization header"})
			return
		}
		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuthentication, token)
		if err != nil {
			utils.WriteJsonResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid token"})
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)

	})
}

func (um *UserMiddleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJsonResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	}
}
