package gc

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PhoneNumber    int       `json:"phoneNumber"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	gorm.DeletedAt `json:"deletedAt"`
}

type ListUsers struct {
	Data  []Users `json:"data"`
	Total int64   `json:"total"`
}

func (g *GC) NewUser(payload *RequestUser) (userId string, err error) {
	hashPass, err := g.HashPassword(payload.Password)
	if err != nil {
		log.Println(err)
		return
	}

	user := &Users{
		ID:          uuid.NewString(),
		Username:    payload.Username,
		Password:    hashPass,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = g.Db.Create(&user).Error
	if err != nil {
		log.Println(err)
		return "", err
	}

	return user.ID, nil
}

func (g *GC) GetUser(userId string) (user *Users, err error) {
	err = g.Db.First(&user, "id = ?", userId).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (g *GC) UpdateUser(userId string, payload *RequestUser) (user *Users, err error) {
	err = g.Db.First(&user, "id = ?", userId).Error
	if err != nil {
		log.Println(err)
		return
	}

	hashPass, err := g.HashPassword(payload.Password)
	if err != nil {
		log.Println(err)
		return
	}
	user.Username = payload.Username
	user.Password = hashPass
	user.Email = payload.Email
	user.PhoneNumber = payload.PhoneNumber

	err = g.Db.Save(&user).Error
	if err != nil {
		log.Println(err)
		return user, err
	}

	return user, nil
}

func (g *GC) DeleteUser(userId string) (err error) {
	var user Users
	err = g.Db.First(&user, "id = ?", userId).Error
	if err != nil {
		log.Println(err)
		return
	}

	user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	err = g.Db.Save(&user).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (g *GC) ListUser(keyword string, limit, page int) (users ListUsers, err error) {
	err = g.Db.Table("users").Where("deleted_at IS NULL AND username LIKE ?", "%"+keyword+"%").Limit(limit).Offset((page - 1) * limit).Find(&users.Data).Error
	if err != nil {
		log.Println(err)
		return
	}

	err = g.Db.Table("users").Where("deleted_at IS NULL AND username LIKE ?", "%"+keyword+"%").Count(&users.Total).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (g *GC) Auth(payload *AuthUser) (token string, err error) {
	var user Users
	err = g.Db.First(&user, "username = ?", payload.Username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.New("please check your username or password")
		}
		log.Println(err)
		return "", err
	}

	log.Println("user =>", user)

	if !g.CheckPassword(user.Password, payload.Password) {
		err = errors.New("please check your username or password")
		log.Println(err)
		return "", err
	}

	token, err = g.CreateJWT(user.ID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return token, nil

}
