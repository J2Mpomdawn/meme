package control

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

func ReConn(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html.tmpl", gin.H{
		"Region": "re-connect",
	})

	service.DbEngine.Close()
	service.Connect()
}

func Com(c *gin.Context) {
	var cmd model.Cmd
	c.Bind(&cmd)

	switch cmd.App {
	case "dis":
		if len(cmd.Args) > 0 {
			for _, v := range cmd.Args {
				if v == "-d" {
					service.DbEngine.Close()
				} else if v == "-s" {
					service.SshClient.Close()
				}
			}
		} else {
			service.DisConnect()
		}

	case "conn":
		service.Connect()
	}
}
