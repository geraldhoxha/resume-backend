package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/geraldhoxha/resume-backend/service"
)

type authString string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		auth := r.Header.Get("Authorization")

		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		bearer := "Bearer "

		auth = auth[len(bearer):]
		validate, err := service.JwtValidate(context.Background(), auth)
		if err != nil || !validate.Valid{
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		customClaim, ok := validate.Claims.(*service.JwtCustomClaim)
		if !ok {
			http.Error(w, "Something went wrong. Please login", http.StatusNotFound)
			return 
		}
		ctx := context.WithValue(r.Context(), authString("auth"), customClaim)
		r = r.WithContext(ctx)
		next.ServeHTTP(w,r)
	})
}

func CtxValue(ctx context.Context) *service.JwtCustomClaim {
	raw, ok := ctx.Value(authString("auth")).(*service.JwtCustomClaim)
	if !ok {
		fmt.Println("Error bro", ctx.Value(authString("auth")))
	}
	return raw
}
