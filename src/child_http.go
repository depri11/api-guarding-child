package gc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (g *GC) RegisterChildRouter(publicApiRouter, protectedApiRouter *mux.Router) {
	protectedApiRouter.Use(g.AuthMiddleware)

	protectedApiRouter.Path("/child/new").Methods("POST").HandlerFunc(g.CreateChildHandler)
	protectedApiRouter.Path("/child/get/{id}").Methods("GET").HandlerFunc(g.GetChildHandler)
	protectedApiRouter.Path("/child/update/{id}").Methods("PUT").HandlerFunc(g.UpdateChildHandler)
	protectedApiRouter.Path("/child/delete/{id}").Methods("DELETE").HandlerFunc(g.DeleteChildHandler)
	protectedApiRouter.Path("/child/update/me").Methods("PUT").HandlerFunc(g.UpdateChildHandler)
	protectedApiRouter.Path("/child/me").Methods("GET").HandlerFunc(g.MyChildsHandler)
}

type RequestChild struct {
	PhoneNumber string `json:"phoneNumber"`
	ParentId    string `json:"parentId"`
}

func (g *GC) CreateChildHandler(w http.ResponseWriter, r *http.Request) {
	parentId, ok := r.Context().Value("userId").(string)
	if !ok {
		sendGenericHTTPError(w, http.StatusInternalServerError, fmt.Errorf("cannot get userId from request"))
		return
	}

	var payload RequestChild
	bodyBytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	payload.ParentId = parentId

	var childId string

	child, err := g.GetChild(payload.PhoneNumber)
	if err == gorm.ErrRecordNotFound {
		childId, err = g.NewChild(&payload)
		if err != nil {
			sendGenericHTTPError(w, http.StatusInternalServerError, err)
			return
		}
		sendGenericHTTPOk(w, childId)
		return
	}

	b, err := json.Marshal(child)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) GetChildHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	child, err := g.GetChild(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			sendGenericHTTPError(w, http.StatusNotFound, err)
			return
		}
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(child)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) UpdateChildHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var payload RequestChild
	bodyByte, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(bodyByte, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	child, err := g.UpdateChild(id, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(child)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) DeleteChildHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := g.DeleteChild(id)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	sendGenericHTTPOk(w, "ok!")
}

func (g *GC) MyChildsHandler(w http.ResponseWriter, r *http.Request) {
	parentId, ok := r.Context().Value("userId").(string)
	if !ok {
		sendGenericHTTPError(w, http.StatusInternalServerError, fmt.Errorf("cannot get userId from request"))
		return
	}

	childs, err := g.MyChilds(parentId, 0, 0)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(childs)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)

}
