package gc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (g *GC) RegisterStatussRouter(publicApiRouter, protectedApiRouter *mux.Router) {
	publicApiRouter.Path("/status").Methods("GET").HandlerFunc(g.StatusHandler)
	publicApiRouter.Path("/user").Methods("POST").HandlerFunc(g.CreateUserHandler)
}

type RequestRegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *GC) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mapMsg := map[string]interface{}{
		"status":  200,
		"message": "Hello, world!",
	}

	b, err := json.Marshal(mapMsg)
	if err != nil {
		return
	}
	w.Write(b)
}

func (g *GC) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestRegisterUser
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	userId, err := g.NewUser(&payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	sendGenericHTTPOk(w, userId)
}
