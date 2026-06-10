package response

import "github.com/gin-gonic/gin"

func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}