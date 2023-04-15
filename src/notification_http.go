package gc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (g *GC) RegisterNotificationRouter(publicApiRouter, protectedApiRouter *mux.Router) {
	protectedApiRouter.Use(g.AuthMiddleware)

	protectedApiRouter.Path("/notification/new").Methods("POST").HandlerFunc(g.CreateNotificationHandler)
	protectedApiRouter.Path("/notification/{id}").Methods("GET").HandlerFunc(g.ListNotificationByChildHandler)
}

func (g *GC) CreateNotificationHandler(w http.ResponseWriter, r *http.Request) {
	parentId, ok := r.Context().Value("userId").(string)
	if !ok {
		sendGenericHTTPError(w, http.StatusInternalServerError, fmt.Errorf("cannot get userId from request"))
		return
	}

	var payload Notification
	bodyBytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	payload.ID = uuid.NewString()
	payload.ParentId = parentId

	notifId, err := g.NewNotification(&payload)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	sendGenericHTTPOk(w, notifId)
}

func (g *GC) ListNotificationByChildHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	listNotif, err := g.ListNotificationByChild(id)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(listNotif)
	if err != nil {
		sendGenericHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	httpWrite(w, b)
}
