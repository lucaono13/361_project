package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lucaono13/361_project/handlers"
)

type Exception struct {
	Message string `json:"message"`
}

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("x-access-token")

		header = strings.TrimSpace(header)

		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
			return
		}

		tok := &structure.SignedDetails{}

		_, err := jwt.ParseWithClaims(header, tok, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tok)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
