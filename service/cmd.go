package service

import (
	"os/exec"

	"meme/model"
)

// execute command
func ExecCmd(cmd model.Cmd) ([]byte, error) {
	//execute received parameter
	out, err := exec_cmd(cmd)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// error message and proposed amendment for app
func ValidAppMessage(app string) []string {
	//check if app is registered
	if _, ok := model.CmdList[app]; ok {
		return nil
	}

	//list app names
	apps := make([]string, len(model.CmdList))
	{
		i := 0
		for k := range model.CmdList {
			apps[i] = k
			i++
		}
	}

	//comparing app and reserved words
	appearance_count := AppearanceCount(app, apps...)

	//sort candidates
	return valid_message(appearance_count, app)
}

// error message and proposed amendment for argument
func ValidArgMessage(app string, arg []string) []string {
	//list app args
	args := make(map[string]struct{}, len(model.CmdList[app]))
	for _, v := range model.CmdList[app] {
		args[v] = struct{}{}
	}

	for _, v := range arg {
		//check if arg is registered
		if _, ok := args[v]; ok {
			continue
		}

		//comparing argument and reserved words
		appearance_count := AppearanceCount(v, model.CmdList[app]...)

		return valid_message(appearance_count, v)
	}
	return nil
}

/*
func appearance_count(arg string) map[string]int {
	country := map[rune]struct{}{45: {}, 99: {}, 110: {}, 111: {}, 116: {}, 117: {}, 121: {}}
	world := map[rune]struct{}{45: {}, 100: {}, 108: {}, 111: {}, 114: {}, 119: {}}
	group := map[rune]struct{}{45: {}, 103: {}, 111: {}, 112: {}, 114: {}, 117: {}}
	class := map[rune]struct{}{45: {}, 97: {}, 99: {}, 108: {}, 115: {}}
	block := map[rune]struct{}{45: {}, 98: {}, 99: {}, 107: {}, 108: {}, 111: {}}
	castle := map[rune]struct{}{45: {}, 97: {}, 99: {}, 101: {}, 108: {}, 115: {}, 116: {}}

	appearance_count := map[string]int{}
	appearance_count["country"] = 0
	appearance_count["world"] = 0
	appearance_count["group"] = 0
	appearance_count["class"] = 0
	appearance_count["block"] = 0
	appearance_count["castle"] = 0

	for _, b := range arg {
		_, ok := country[b]
		if ok {
			appearance_count["country"]++
		}
		_, ok = world[b]
		if ok {
			appearance_count["world"]++
		}
		_, ok = group[b]
		if ok {
			appearance_count["group"]++
		}
		_, ok = class[b]
		if ok {
			appearance_count["class"]++
		}
		_, ok = block[b]
		if ok {
			appearance_count["block"]++
		}
		_, ok = castle[b]
		if ok {
			appearance_count["castle"]++
		}
	}

	return appearance_count
}
*/
// Command Execution Core
func exec_cmd(cmd model.Cmd) ([]byte, error) {
	//command execution
	out, err := exec.Command(cmd.App, cmd.Args...).Output()

	return out, err
}

// message
func valid_message(appearance_count map[string]int, arg string) []string {
	//sort candidates
	pl := SortMapValue_StrInt(appearance_count)
	i := len(model.CmdList) - 1

	FmtPrintln("bluee", StrJoin(128, "\"", arg, "\"", " is not a valid argument\nmaybe: \"", pl[i].Key, "\""))

	return []string{
		StrJoin(32, "\"", arg, "\"", " is not a valid argument"),
		StrJoin(32, "maybe: \"", pl[i].Key, "\""),
	}
}
