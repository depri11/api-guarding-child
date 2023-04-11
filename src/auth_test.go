package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	db, err := test_initDb()
	assert.Equal(t, nil, err)

	newUser := &RequestUser{
		Username:    "test",
		Password:    "test123",
		Email:       "test@gmail.com",
		PhoneNumber: "62822",
	}

	svc := &GC{
		Db: db,
	}

	userId, err := svc.NewUser(newUser)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", userId)

	payload := &AuthUser{
		Username: "test",
		Password: "test123",
	}

	token, err := svc.Auth(payload)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", token)
}

func TestHashPassword(t *testing.T) {
	svc := &GC{}

	myPassword := "Test Password"

	hash, err := svc.HashPassword(myPassword)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", hash)

	ok := svc.CheckPassword(hash, "Test Password")
	assert.Equal(t, true, ok)
}
