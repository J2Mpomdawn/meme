package control

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

// processing of POST request for data communication
func Com(c *gin.Context) {
	//bind received command arguments to model
	var cmd model.Cmd
	c.Bind(&cmd)

	switch cmd.App {
	case "dis":
		//disconnect
		if len(cmd.Args) > 0 {
			for _, v := range cmd.Args {
				if v == "-d" {
					//disconnect DB
					service.DbEngine.Close()

				} else if v == "-s" {
					//disconnect SSH
					service.SshClient.Close()
				}
			}
		} else {
			//Disconnect Both
			service.DisConnect()
		}

	case "conn":
		//Connect DB and SSH
		service.Connect()
	}
}

// processing of POST request for reconnection
func ReConn(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html.tmpl", gin.H{
		"Region": "re-connect",
	})

	service.DbEngine.Close()
	service.Connect()
}
