package gc

import (
	"log"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type GC struct {
	Db     *gorm.DB
	Router *mux.Router
	Cache  *cache.Cache
}

func NewGC(db *gorm.DB) *GC {
	err := db.AutoMigrate(&Users{}, &Child{}, &Location{}, &Notification{})
	if err != nil {
		log.Println(err)
	}
	cache := cache.New(24*time.Hour, 1*time.Hour)
	g := &GC{
		Db:     db,
		Router: mux.NewRouter(),
		Cache:  cache,
	}

	g.RegisterRouter()

	return g
}
