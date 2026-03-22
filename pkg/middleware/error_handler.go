package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*AppError); ok {
				c.JSON(appErr.HTTPStatus, gin.H{
					"success": false,
					"error":   appErr,
				})
				return
			}

			// fallback
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    500,
					"message": err.Error(),
				},
			})
		}
	}
}
