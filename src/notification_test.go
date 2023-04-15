package gc

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotificationService(t *testing.T) {
	db, err := test_initDb()
	assert.Equal(t, nil, err)

	svc := &GC{
		Db: db,
	}

	notif := &Notification{
		ID:          uuid.NewString(),
		ParentId:    uuid.NewString(),
		ChildId:     uuid.NewString(),
		Title:       "Notif title",
		Text:        "Description notif",
		PackageName: "com.example.aplication",
		Timestamp:   time.Now().Unix(),
	}

	id, err := svc.NewNotification(notif)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, id)

	notif.ID = uuid.NewString()
	notif.Title = "Notif title 2"
	notif.Text = "Description notif 2"
	notif.PackageName = "com.example.aplication2"

	id2, err := svc.NewNotification(notif)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, id2)

	listNotif, err := svc.ListNotificationByChild(notif.ChildId)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(listNotif.Data))
	assert.Equal(t, int64(2), listNotif.Total)

}
