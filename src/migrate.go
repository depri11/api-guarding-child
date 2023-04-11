package gc

import (
	"log"

	"gorm.io/gorm"
)

func dbMigrate(db *gorm.DB) (err error) {
	// auto migrate
	err = db.Migrator().
		DropTable(&Users{}, &Child{})
	if err != nil {
		log.Println(err)
	}
	err = db.AutoMigrate(Users{}, &Child{})
	if err != nil {
		log.Println(err)
	}
	return
}
