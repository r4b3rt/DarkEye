package main

import (
	"github.com/c-bata/go-prompt"
)

var (
	analysisProgram     = "analysis"
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

func (ctx *RequestContext) analysisArgumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return filterSuggestions(runCompleteCheck(analysisSuggestions, ctx.CmdArgs,
			[]string{
				"-sql",
			}), ctx.CmdArgs)
	}
	switch args[0] {
	case "-sql":
		if len(args) == 2 {
			return []prompt.Suggest{{"select * from ent limit 1", "Analysis the result of intelligence collection"},}
		}
	}
	return []prompt.Suggest{}
}
