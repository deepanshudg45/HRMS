// internal/auth/handler_test.go
package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	h := Handler{/* mock repo */}

	router.POST("/login", h.Login)

	body := `{"work_email":"test@test.com","password":"123456"}`

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}