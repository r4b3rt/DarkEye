package main

import (
	"github.com/c-bata/go-prompt"
)

var (
	analysisSuggestions = []prompt.Suggest{
		// Command
		{"-sql", "select * from ent limit 1"},
		{"-output-csv", "output.csv"},
	}
	analysisValueCheck = map[string]bool{
		"-sql":        false,
		"-output-csv": false,
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
	if isDuplicateArg(args[0], a.parent.CmdArgs) {
		return []prompt.Suggest{}
	}
	switch args[0] {
	case "-sql":
		if len(args) == 2 {
			return []prompt.Suggest{
				{"select * from ent limit 1", "SQL语句查询并输出结果(支持ent"},
				{"delete from ent", "删掉所有数据"},}
		}
	case "-output-csv":
		if len(args) == 2 {
			return []prompt.Suggest{
				{"output.csv", "输出查询结果"},
			}
		}

	}
	return []prompt.Suggest{}
}
