package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
)

type AuthMiddleware struct {
	SignKey interface{}
}

func (mw AuthMiddleware) Authenticated(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return mw.SignKey, nil
		})

		if err != nil {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			ApiResponse(w, "You are not allowed", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		p = append(p, httprouter.Param{Key: "username", Value: claims["username"].(string)})
		p = append(p, httprouter.Param{Key: "userId", Value: strconv.Itoa(int(claims["userId"].(float64)))})
		handle(w, r, p)
	}
}
