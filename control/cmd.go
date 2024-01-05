package control

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

var TestCounter = 0

func Cmd(c *gin.Context) {
	var cmd model.Cmd
	c.Bind(&cmd)

	var res gin.H

	switch cmd.App {
	case "help":
		switch c.Request.URL.Path {
		case service.StrJoin(32, "/", os.Getenv("Master")):
			res = help_mst(cmd)

		case "/cmd/exec":
			res = help_cmd(cmd)
		}

	case "test":
		switch cmd.Args[0] {
		case "-s":
			TestCounter++
		case "-c":
			res = gin.H{
				"test-counter": TestCounter,
			}
		}

	default:
		res = exec(cmd)
	}

	c.JSON(http.StatusOK, res)
}

func exec(cmd model.Cmd) gin.H {
	service.ExecCmd(cmd)

	return gin.H{
		"test_res":  "aiueo",
		"test_res2": "oooo",
	}
}

func help_cmd(cmd model.Cmd) gin.H {
	return gin.H{
		"cmd1": "cmd1-desc",
		"cmd2": "cmd2-desc",
		"help": "Provides Help information for key commands.",
		"etc":  "...",
	}
}

func help_mst(cmd model.Cmd) gin.H {
	return gin.H{
		"e":    "Get",
		"ex":   "out",
		"exi":  "of",
		"exit": "room.",
	}
}
