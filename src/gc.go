package gc

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GC struct {
	Db       *gorm.DB
	Router   *mux.Router
	Fixtures Fixtures
	Cache    *cache.Cache
}

func NewGC(db *gorm.DB) *GC {
	// usersService := svc.NewService(db)
	cache := cache.New(24*time.Hour, 1*time.Hour)
	g := &GC{
		Db:     db,
		Router: mux.NewRouter(),
		Cache:  cache,
	}

	g.RegisterRouter()

	return g
}

func InitDB() (*gorm.DB, error) {
	dbDriver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		return nil, errors.New("missing-dsn")
	}

	var dialector gorm.Dialector
	switch dbDriver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		log.Println(dbDriver)
		return nil, errors.New("unknown-driver")
	}

	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, err
}
