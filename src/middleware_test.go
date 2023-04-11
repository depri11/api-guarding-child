package gc

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	db, err := test_initDb()
	assert.Equal(t, nil, err)

	svc := &GC{
		Db: db,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("entered handler")
		ctx := r.Context()
		assert.NotNil(t, ctx)
		userid := r.Context().Value("userId")
		assert.NotEmpty(t, userid)
	})
	handlerTest := svc.AuthMiddleware(handler)

	// test with no auth header
	req := httptest.NewRequest("GET", "http://testing", nil)
	rec := httptest.NewRecorder()
	handlerTest.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Result().StatusCode)

	// test with invalid auth header
	req = httptest.NewRequest("GET", "http://testing", nil)
	rec = httptest.NewRecorder()
	req.Header.Add("Authorization", "Bearer fake")
	handlerTest.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Result().StatusCode)

	token, err := svc.CreateJWT("userID")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", token)
	req = httptest.NewRequest("GET", "http://testing", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	handlerTest.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
