package model

//used for command execution
type Cmd struct {
	//app name
	App string `app:"app"`
	//arguments
	Args []string `args:"args"`
}

//command list
var CmdList = map[string][]string{
	"help":   {},
	"check":  {},
	"select": {},
	"...":    {},
}
