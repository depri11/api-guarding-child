package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceUser(t *testing.T) {
	db, err := test_initDb()
	assert.Equal(t, nil, err)

	newUser := &RequestUser{
		Username:    "test",
		Password:    "test123",
		Email:       "test@gmail.com",
		PhoneNumber: 62822,
	}

	svc := &GC{
		Db: db,
	}

	userId, err := svc.NewUser(newUser)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", userId)

	user, err := svc.GetUser(userId)
	assert.Equal(t, nil, err)
	assert.Equal(t, newUser.Username, user.Username)
	assert.NotEqual(t, newUser.Password, user.Password)
	assert.NotEqual(t, "", user.Password)
	assert.Equal(t, newUser.Email, user.Email)
	assert.Equal(t, newUser.PhoneNumber, user.PhoneNumber)

	updateUser := &RequestUser{
		Username: "updateTest",
	}

	user, err = svc.UpdateUser(userId, updateUser)
	assert.Equal(t, nil, err)
	assert.Equal(t, updateUser.Username, user.Username)

	users, err := svc.ListUser("", 0, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, user.Username, users.Data[0].Username)
	assert.Equal(t, user.Password, users.Data[0].Password)
	assert.Equal(t, user.Email, users.Data[0].Email)
	assert.Equal(t, user.PhoneNumber, users.Data[0].PhoneNumber)
	assert.Equal(t, int64(1), users.Total)

	err = svc.DeleteUser(userId)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, user.DeletedAt)

}
