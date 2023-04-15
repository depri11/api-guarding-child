package gc

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Child struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	ParentID    string         `json:"parentId"`
	PhoneNumber string         `json:"phoneNumber"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}

type ListMyChilds struct {
	Data  []Child `json:"data"`
	Total int64   `json:"total"`
}

func (g *GC) NewChild(payload *RequestChild) (userId string, err error) {
	child := &Child{
		ID:          uuid.NewString(),
		PhoneNumber: payload.PhoneNumber,
		ParentID:    payload.ParentId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = g.Db.Create(&child).Error
	if err != nil {
		log.Println(err)
		return "", err
	}

	return child.ID, nil
}

func (g *GC) GetChild(childId string) (child *Child, err error) {
	err = g.Db.First(&child, "id = ? OR phone_number = ?", childId, childId).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (g *GC) UpdateChild(childId string, payload *RequestChild) (child *Child, err error) {
	err = g.Db.First(&child, "id = ?", childId).Error
	if err != nil {
		log.Println(err)
		return
	}

	child.PhoneNumber = payload.PhoneNumber

	err = g.Db.Save(&child).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (g *GC) DeleteChild(userId string) (err error) {
	var child Child
	err = g.Db.First(&child, "id = ?", userId).Error
	if err != nil {
		log.Println(err)
		return
	}

	child.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	err = g.Db.Save(&child).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (g *GC) MyChilds(userId string, limit, page int) (childs ListMyChilds, err error) {
	err = g.Db.Model(&Child{}).Where("parent_id = ?", userId).Limit(limit).Offset((page - 1) * limit).Find(&childs.Data).Error
	if err != nil {
		log.Println(err)
		return
	}

	err = g.Db.Model(&Child{}).Where("deleted_at IS NULL AND parent_id = ?", userId).Count(&childs.Total).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}
