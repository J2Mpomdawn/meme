package control

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// processing of GET request for index page
func Index(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
		"Path":   path,
		"Region": "index",
	})
}
