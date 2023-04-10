package gc

import (
	"errors"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
