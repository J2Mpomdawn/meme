package control

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Key(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "key.html.tmpl", gin.H{
		"Path":   path,
		"Region": "key",
	})

	//service.SSHPermission()
}
