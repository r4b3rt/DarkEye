package main

import (
	"github.com/c-bata/go-prompt"
)

var (
	analysisSuggestions = []prompt.Suggest{
		// Command
		{"-sql", "select * from ent limit 1"},
	}
	analysisValueCheck = map[string]bool{
		"-sql": false,
	}
)

func init() {
	mSuggestions = append(mSuggestions, []prompt.Suggest{
		{analysisProgram, "Analysis the result of intelligence collection"},
	}...)
}

func (a *analysisRuntime) Completer(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return filterSuggestions(runCompleteCheck(analysisSuggestions, a.parent.CmdArgs,
			[]string{
				"-sql",
			}), a.parent.CmdArgs)
	}
	//过滤重复的命令
	if len(filterSuggestions([]prompt.Suggest{
		{args[0],"any"},
	}, a.parent.CmdArgs)) == 0 {
		return []prompt.Suggest{}
	}
	switch args[0] {
	case "-sql":
		if len(args) == 2 {
			return []prompt.Suggest{{"select * from ent limit 1", "Analysis the result of intelligence collection"},}
		}
	}
	return []prompt.Suggest{}
}
