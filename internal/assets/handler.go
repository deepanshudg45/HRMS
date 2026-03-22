package assets

import "github.com/gin-gonic/gin"

type Handler struct {
	Service *Service
}

func (h *Handler) GetAssets(c *gin.Context) {
	c.JSON(200, gin.H{"message": "assets list"})
}
