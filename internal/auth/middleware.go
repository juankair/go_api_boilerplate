package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/juankair/go_api_boilerplate/pkg/response"
	"github.com/uptrace/bunrouter"
	"net/http"
	"strings"
)

func SecureMiddleware(jwtKey string) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			authHeader := req.Header.Get("authorization")
			if authHeader == "" {
				response.RespondWithJSON(w, http.StatusUnauthorized, false, "Authorization Invalid", nil)
				return nil
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 {
				response.RespondWithJSON(w, http.StatusUnauthorized, false, "Authorization Invalid", nil)
				return nil
			}

			tokenStr := bearerToken[1]
			claims := &jwt.StandardClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})

			if err != nil {
				response.RespondWithJSON(w, http.StatusUnauthorized, false, "Authorization Invalid", nil)
				return nil
			}

			if !token.Valid {
				response.RespondWithJSON(w, http.StatusUnauthorized, false, "Token Expired", nil)
				return nil
			}

			if claims.Id == "" {
				response.RespondWithJSON(w, http.StatusUnauthorized, false, "AccountID is empty", nil)
				return nil
			}

			ctx := context.WithValue(req.Context(), "accountID", claims.Id)
			return next(w, req.WithContext(ctx))
		}
	}
}
