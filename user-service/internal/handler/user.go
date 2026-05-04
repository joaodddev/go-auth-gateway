package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users := []map[string]interface{}{
		{"id": 1, "name": "João"},
		{"id": 2, "name": "Maria"},
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}