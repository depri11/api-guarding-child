package gc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (g *GC) RegisterStatussRouter(publicApiRouter, protectedApiRouter *mux.Router) {
	publicApiRouter.Path("/status").Methods("GET").HandlerFunc(g.StatusHandler)
	publicApiRouter.Path("/user/new").Methods("POST").HandlerFunc(g.CreateUserHandler)
	publicApiRouter.Path("/user/auth").Methods("POST").HandlerFunc(g.AuthHandler)

	protectedApiRouter.Use(g.AuthMiddleware)

	protectedApiRouter.Path("/users/get/{id}").Methods("GET").HandlerFunc(g.GetUserHandler)
	protectedApiRouter.Path("/users/update/{id}").Methods("PUT").HandlerFunc(g.UpdateUserHandler)
	protectedApiRouter.Path("/users/delete/{id}").Methods("DELETE").HandlerFunc(g.DeleteUserHandler)
	protectedApiRouter.Path("/users/update/me").Methods("PUT").HandlerFunc(g.UpdateUserHandler)
	protectedApiRouter.Path("/users/list").Methods("PUT").HandlerFunc(g.ListUserHandler)
}

type RequestUser struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber int    `json:"phoneNumber"`
}

type AuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *GC) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mapMsg := map[string]interface{}{
		"message": "Hello, world!",
	}

	b, err := json.Marshal(mapMsg)
	if err != nil {
		return
	}
	w.Write(b)
}

func (g *GC) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestUser
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

func (g *GC) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		sendGenericHTTPError(w, http.StatusInternalServerError, fmt.Errorf("cannot get userId from request"))
		return
	}

	id := mux.Vars(r)["id"]
	if id == "me" {
		id = userId
	}

	user, err := g.GetUser(id)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		sendGenericHTTPError(w, http.StatusInternalServerError, fmt.Errorf("cannot get userId from request"))
		return
	}

	id := mux.Vars(r)["id"]
	if id == "me" {
		id = userId
	}

	var payload RequestUser
	bodyByte, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(bodyByte, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := g.UpdateUser(id, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := g.DeleteUser(id)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	sendGenericHTTPOk(w, "ok!")
}

func (g *GC) ListUserHandler(w http.ResponseWriter, r *http.Request) {
	users, err := g.ListUser("", 0, 0)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(users)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var payload AuthUser
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := g.Auth(&payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusForbidden, err)
		return
	}

	sendGenericHTTPOk(w, token)

}
