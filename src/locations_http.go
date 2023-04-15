package gc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (g *GC) RegisterLocationsRouter(publicApiRouter, protectedApiRouter *mux.Router) {
	protectedApiRouter.Use(g.AuthMiddleware)
	protectedApiRouter.Path("/location/new").Methods("POST").HandlerFunc(g.CreateLocationHandler)
	protectedApiRouter.Path("/location/get/{id}").Methods("GET").HandlerFunc(g.GetLocationHandler)
	protectedApiRouter.Path("/location/history/{id}").Methods("GET").HandlerFunc(g.GetHistoryLocationHandler)
	protectedApiRouter.Path("/location/delete/{id}").Methods("DELETE").HandlerFunc(g.DeleteLocationHandler)
}

type RequestLocation struct {
	ParentID string  `json:"parentId"`
	ChildID  string  `json:"childId"`
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
	Address  string  `json:"address"`
}

func (g *GC) CreateLocationHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		sendGenericHTTPError(w, http.StatusInternalServerError, fmt.Errorf("cannot get userId from request"))
		return
	}

	var payload RequestLocation
	bodyBytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	payload.ParentID = userId

	id, err := g.NewLocation(&payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	sendGenericHTTPOk(w, id)
}

func (g *GC) GetLocationHandler(w http.ResponseWriter, r *http.Request) {
	childId := mux.Vars(r)["id"]

	loc, err := g.GetLocation(childId)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(loc)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) GetHistoryLocationHandler(w http.ResponseWriter, r *http.Request) {
	childId := mux.Vars(r)["id"]

	loc, err := g.GetHistoryLocation(childId, 0, 0)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(loc)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}

func (g *GC) DeleteLocationHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := g.DeleteLocation(id)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	sendGenericHTTPOk(w, "ok!")
}
