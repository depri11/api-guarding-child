package gc

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	ParentID  string         `json:"parentId"`
	ChildID   string         `json:"childId"`
	Lat       float64        `json:"lat"`
	Long      float64        `json:"long"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

type HistoryLocation struct {
	Data  []Location `json:"data"`
	Total int64      `json:"total"`
}

func (g *GC) NewLocation(payload *RequestLocation) (id string, err error) {
	location := &Location{
		ID:        uuid.NewString(),
		ChildID:   payload.ChildID,
		ParentID:  payload.ParentID,
		Lat:       payload.Lat,
		Long:      payload.Long,
		Address:   payload.Address,
		CreatedAt: time.Now(),
	}

	err = g.Db.Create(&location).Error
	if err != nil {
		log.Println(err)
		return "", err
	}

	return location.ID, nil
}

func (g *GC) GetLocation(childId string) (loc Location, err error) {
	err = g.Db.Where("child_id = ?", childId).First(&loc).Error
	if err != nil {
		log.Println(err)
		return loc, err
	}

	return loc, nil
}

func (g *GC) GetHistoryLocation(childId string, page, limit int) (history HistoryLocation, err error) {
	err = g.Db.Model(&Location{}).Where("child_id = ?", childId).Limit(limit).Offset((page - 1) * limit).Find(&history.Data).Error
	if err != nil {
		log.Println(err)
		return history, err
	}

	err = g.Db.Model(&Location{}).Where("child_id = ?", childId).Count(&history.Total).Error
	if err != nil {
		log.Println(err)
		return history, err
	}

	return history, nil
}

func (g *GC) DeleteLocation(id string) (err error) {
	var location Location
	err = g.Db.First(&location, "id = ?", id).Error
	if err != nil {
		log.Println(err)
		return
	}

	location.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	err = g.Db.Save(&location).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}
