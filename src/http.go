package gc

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GenericHTTPMessage struct {
	Message      string `json:"message"`
	IsSuccessful bool   `json:"isSuccessful"`
	Code         int    `json:"code"`
	Info         string `json:"info"`
}

func sendGenericHTTPError(w http.ResponseWriter, status int, err error) {
	logrus.Error(err)
	w.WriteHeader(status)

	g := GenericHTTPMessage{
		Message:      err.Error(),
		IsSuccessful: false,
		Code:         status,
	}

	b, _ := json.Marshal(&g)
	httpWrite(w, b)
}

func sendGenericHTTPOk(w http.ResponseWriter, msg string) {
	g := GenericHTTPMessage{
		Message:      msg,
		IsSuccessful: true,
		Code:         200,
	}

	b, err := json.Marshal(&g)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, errors.New("generic-error"))
		return
	}
	httpWrite(w, b)

}

func httpWrite(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (g *GC) RegisterRouter() {
	publicApiRouter := g.Router.PathPrefix("/api/v1/public").Subrouter()
	protectedApiRouter := g.Router.PathPrefix("/api/v1").Subrouter()

	g.RegisterStatussRouter(publicApiRouter, protectedApiRouter)
}
