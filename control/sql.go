package control

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

// execute CRUD
func CRUD(c *gin.Context) {
	//bind received command arguments to model
	var cmd model.Cmd
	c.Bind(&cmd)

	switch cmd.App {
	case "select":
		res, err := service.Select(&cmd.Args)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"failed": [1]string{err.Error()},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"table": res,
			})
		}
	case "insert", "update", "delete":
		//even though it is hidden, it is dangerouus to update data
		//from an environment where anyone can touch it
	}
}

// Deprecated: to be deleted
func SqlPage(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
		"Path":   path,
		"Region": "sql-test",
	})

	service.SelectTest()
}
