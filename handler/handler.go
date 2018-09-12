package handler

import (
	"hamal/models"

	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) *models.User {
	v, ok := c.Get("user")
	if !ok {
		return &models.User{}
	}
	u := v.(*models.User)
	return u
}
