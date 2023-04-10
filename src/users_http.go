package gc

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (g *GC) RegisterStatussRouter(publicApiRouter, protectedApiRouter *mux.Router) {
	publicApiRouter.Path("/status").Methods("GET").HandlerFunc(g.StatusHandler)
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
