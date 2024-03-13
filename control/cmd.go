package control

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

// processing of POST request for command
func Cmd(c *gin.Context) {
	//bind received command arguments to model
	var cmd model.Cmd
	c.Bind(&cmd)

	var res gin.H
	var err_flg bool

	switch cmd.App {
	case "help":
		//how to use it
		switch c.Request.URL.Path {
		case "/cmd/exec":
			res = help_cmd(cmd)
		default:
			res = help_ext(cmd)
		}
	default:
		//oter commands
		res, err_flg = exec(cmd)
	}

	if err_flg {
		c.JSON(http.StatusInternalServerError, res)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

// execute command
func exec(cmd model.Cmd) (gin.H, bool) {
	var res gin.H
	var err_flg bool

	//app validator
	if out := service.ValidAppMessage(cmd.App); out != nil {
		res = gin.H{
			"failed": out,
		}
		//arg validator
	} else if out = service.ValidArgMessage(cmd.App, cmd.Args); out != nil {
		res = gin.H{
			"failed": out,
		}
	} else {
		//execute
		out, err := service.ExecCmd(cmd)

		if err != nil {
			res = gin.H{
				"failed": [...]string{string(out), err.Error()},
			}
			err_flg = true
		} else {
			res = gin.H{
				"": string(out),
			}
		}
	}

	return res, err_flg
}

// cmd help list
func help_cmd(cmd model.Cmd) gin.H {
	return gin.H{
		"cmd1": "cmd1-desc",
		"cmd2": "cmd2-desc",
		"help": "Provides Help information for key commands.",
		"etc":  "...",
	}
}

// default help list
func help_ext(cmd model.Cmd) gin.H {
	return gin.H{
		"e":    "Get",
		"ex":   "out",
		"exi":  "of",
		"exit": "room.",
	}
}
