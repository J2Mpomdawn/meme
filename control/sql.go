package control

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

func SqlPage(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
		"Path":   path,
		"Region": "sql-test",
	})

	service.SelectTest()
}

func CRUD(c *gin.Context) {
	var cmd model.Cmd
	c.Bind(&cmd)

	switch cmd.App {
	case "select":
		service.Select(&cmd.Args)
	}
}
