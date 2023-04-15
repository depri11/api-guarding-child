package gc

import (
	"log"
)

type Notification struct {
	ID          string `json:"id" gorm:"primaryKey"`
	ParentId    string `json:"parentId"`
	ChildId     string `json:"childId"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	PackageName string `json:"packageName"`
	Timestamp   int64  `json:"timestamp"`
}

type ListNotification struct {
	Data  []Notification `json:"data"`
	Total int64          `json:"total"`
}

func (g *GC) NewNotification(payload *Notification) (id string, err error) {
	err = g.Db.Create(&payload).Error
	if err != nil {
		log.Println(err)
		return "", err
	}

	return payload.ID, nil
}

func (g *GC) ListNotificationByChild(childId string) (listNotif ListNotification, err error) {
	err = g.Db.Model(&Notification{}).Where("child_id = ?", childId).Find(&listNotif.Data).Error
	if err != nil {
		log.Println(err)
		return listNotif, err
	}

	err = g.Db.Model(&Notification{}).Where("child_id = ?", childId).Count(&listNotif.Total).Error
	if err != nil {
		log.Println(err)
		return listNotif, err
	}

	return listNotif, nil
}
