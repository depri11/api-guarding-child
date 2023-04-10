package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, err := test_initDb()
	assert.Equal(t, nil, err)

	user := &RequestRegisterUser{
		Username: "test",
		Password: "test123",
	}

	svc := &GC{
		Db: db,
	}

	userId, err := svc.NewUser(user)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", userId)
}
