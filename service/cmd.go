package service

import (
	"os/exec"
	"sort"

	"meme/model"
)

func SortMapValue(targe_map map[string]int) model.PairList {
	pl := make(model.PairList, len(targe_map))
	i := 0

	for k, v := range targe_map {
		pl[i] = model.Pair{k, v}
		i++
	}

	sort.Sort(pl)

	return pl
}

// execute command
func ExecCmd(cmd model.Cmd) {
	//execute received parameter
	out, err := exec_cmd(cmd)

	if err != nil {
		LogPrint("red", "exec_cmd", err)
	}

	FmtPrint("", cmd)
	FmtPrint("", string(out))
}

// Command Execution Core
func exec_cmd(cmd model.Cmd) ([]byte, error) {
	//command execution
	out, err := exec.Command(cmd.App, cmd.Args...).Output()

	return out, err
}
