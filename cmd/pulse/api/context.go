package api

import (
	"context"
	"github.com/mhamm84/pulse-api/internal/data"
	"net/http"
)

type contextUser string

const userContextKey = contextUser("user")

func (app *application) contextSetUser(r *http.Request, u *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, u)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	u, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user in request context")
	}
	return u
}
