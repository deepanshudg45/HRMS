// internal/auth/handler.go
package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wil/hrms/pkg/utils"
)

type Handler struct {
	Repo *Repository
}

// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hash, _ := utils.HashPassword(req.Password)

	err := h.Repo.CreateEmployee(c, req.WorkEmail, hash, req.EmployeeCode)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "registered"})
}

// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	_ = c.ShouldBindJSON(&req)

	user, err := h.Repo.GetByEmail(c, req.WorkEmail)
	if err != nil || !utils.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	access, _ := GenerateAccessToken(user.ID)
	refresh, _ := GenerateRefreshToken(user.ID)

	hash, _ := HashToken(refresh)

	h.Repo.UpdateRefreshToken(c, user.ID, hash, time.Now().Add(7*24*time.Hour).String())

	c.JSON(200, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// POST /auth/logout
func (h *Handler) Logout(c *gin.Context) {
	empID := c.GetString("employee_id")
	h.Repo.ClearRefreshToken(c, empID)

	c.JSON(200, gin.H{"message": "logged out"})
}

func (h *Handler) Refresh(c *gin.Context) {
	token := c.GetHeader("x-refresh-token")

	// fetch from DB + compare hash (pseudo)
	// if valid → issue new access token
}