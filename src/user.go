package gc

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phoneNumber"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	gorm.DeletedAt `json:"deletedAt"`
}

func (g *GC) NewUser(payload *RequestRegisterUser) (userId string, err error) {
	user := &Users{
		ID:        uuid.NewString(),
		Username:  payload.Username,
		Password:  payload.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = g.Db.Create(&user).Error
	if err != nil {
		log.Println(err)
		return "", err
	}

	return payload.Username, nil
}
