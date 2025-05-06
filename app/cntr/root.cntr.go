package cntr

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Page": "home",
	})
}
