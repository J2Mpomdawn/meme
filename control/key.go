package control

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// processing of GET request for key page
func Key(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "key.html.tmpl", gin.H{
		"Path": path,
	})
}
