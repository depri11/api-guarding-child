package gc

import (
	"log"

	"gorm.io/gorm"
)

func dbMigrate(db *gorm.DB) (err error) {
	// auto migrate
	err = db.Migrator().
		DropTable(&Users{})
	if err != nil {
		log.Println(err)
	}
	err = db.AutoMigrate(Users{})
	if err != nil {
		log.Println(err)
	}
	return
}
