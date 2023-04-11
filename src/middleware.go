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
			w.WriteHeader(http.StatusForbidden)
			return
		}
		userId, err := g.ValidateJWT(token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "userId", userId))
		next.ServeHTTP(w, r)
	})
}
