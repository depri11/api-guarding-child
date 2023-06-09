package gc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (g *GC) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		splitted := strings.Split(token, " ")
		if len(splitted) == 2 {
			token = splitted[1]
		} else {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			err := fmt.Errorf("invalid token")
			log.Println(err)
			sendGenericHTTPError(w, http.StatusForbidden, err)
			return
		}
		claims, err := g.ValidateJWT(token)
		if err != nil {
			sendGenericHTTPError(w, http.StatusForbidden, err)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "userId", claims.UserID))
		next.ServeHTTP(w, r)
	})
}
