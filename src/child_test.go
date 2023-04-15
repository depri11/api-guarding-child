package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChildService(t *testing.T) {
	db, err := test_initDb()
	assert.Equal(t, nil, err)

	newChild := &RequestChild{
		PhoneNumber: "62822",
		ParentId:    "parent id",
	}

	svc := &GC{
		Db: db,
	}

	childId, err := svc.NewChild(newChild)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", childId)

	child, err := svc.GetChild(childId)
	assert.Equal(t, nil, err)
	assert.Equal(t, newChild.PhoneNumber, child.PhoneNumber)
	assert.Equal(t, newChild.ParentId, child.ParentID)

	child, err = svc.GetChild(newChild.PhoneNumber)
	assert.Equal(t, nil, err)
	assert.Equal(t, newChild.PhoneNumber, child.PhoneNumber)
	assert.Equal(t, newChild.ParentId, child.ParentID)

	// updateChild := &RequestChild{
	// 	PhoneNumber: "1 update",
	// 	ParentId:    "parent id",
	// }

	// child, err = svc.UpdateChild(childId, updateChild)
	// assert.Equal(t, nil, err)
	// assert.Equal(t, updateChild.PhoneNumber, child.PhoneNumber)

	// childs, err := svc.MyChilds(newChild.ParentId, 0, 0)
	// assert.Equal(t, nil, err)
	// assert.Equal(t, childId, childs.Data[0].ID)
	// assert.Equal(t, updateChild.ParentId, childs.Data[0].ParentID)
	// assert.Equal(t, updateChild.PhoneNumber, childs.Data[0].PhoneNumber)
	// assert.Equal(t, int64(1), childs.Total)

	// err = svc.DeleteChild(childId)
	// assert.Equal(t, nil, err)
	// assert.NotEqual(t, nil, child.DeletedAt)

}
