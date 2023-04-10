package gc

import (
	"log"
	"os"

	"gorm.io/gorm"
)

func test_initDb() (*gorm.DB, error) {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	os.Setenv("DB_DRIVER", "sqlite")
	os.Setenv("DB_DSN", "file::memory:?cache=shared")
	db, err := InitDB()
	if err != nil {
		return nil, err
	}

	err = dbMigrate(db)
	return db, err
}
