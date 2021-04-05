package main

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

var mSuggestions = []prompt.Suggest{
	{"exit", "Exit"},
	{"cd", "Backward or forward"},
	{"stop", "Stop the running module"},
}

func (ctx *RequestContext) completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	if len(ctx.CmdArgs) > 0 {
		//子命令
		return filterHasPrefix(
			ModuleFuncs[moduleId(ctx.CmdArgs[0])].completer(args), d.GetWordBeforeCursor(),
			moduleId(ctx.CmdArgs[0]))
	}
	return filterHasPrefix(
		mSuggestions, d.GetWordBeforeCursor(), "")
}

//filterHasPrefix 筛选子命令
func filterHasPrefix(suggestions []prompt.Suggest, sub string, moduleId string) []prompt.Suggest {
	if sub == "?" || sub == " " {
		//此时无任何子命令
		return suggestions
	}
	ret := make([]prompt.Suggest, 0)
	for i := range suggestions {
		a := strings.ToLower(suggestions[i].Text)
		b := strings.ToLower(sub)
		if strings.HasPrefix(a, b) {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}

//filterSuggestions 过滤使用过的子命令
func filterSuggestions(suggestions []prompt.Suggest, sub []string) []prompt.Suggest {
	if sub == nil || len(sub) == 1 {
		//此时无任何子命令
		return suggestions
	}
	ok := true
	ret := make([]prompt.Suggest, 0)
	for i := range suggestions {
		ok = true
		a := strings.ToLower(suggestions[i].Text)
		for j := range sub {
			b := strings.ToLower(strings.Split(sub[j], " ")[0])
			if strings.Compare(a, b) == 0 {
				ok = false
			}
		}
		if ok {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}

//检查是否可以结束运行
func runCompleteCheck(suggestions []prompt.Suggest, cmdArgs []string, requirement []string) []prompt.Suggest {
	ret := make([]prompt.Suggest, 0)
	ret = append(ret, suggestions...)
	if len(requirement) == 0 {
		ret = append(ret, prompt.Suggest{
			Text: "exploit", Description: "Do the exploit",
		})
		return ret
	}
	ok := false
	for i := range requirement {
		a := requirement[i]
		ok = false
		for j := range cmdArgs {
			b := strings.ToLower(strings.Split(cmdArgs[j], " ")[0])
			if strings.Compare(a, b) == 0 {
				ok = true
				break
			}
		}
		if !ok {
			return ret
		}
	}
	ret = append(ret, prompt.Suggest{
		Text: "exploit", Description: "Do the exploit",
	})
	return ret
}

func (ctx *RequestContext) livePrefix() (string, bool) {
	if len(ctx.CmdArgs) == 0 {
		return "", false
	}
	promptStr := ""
	for _, v := range ctx.CmdArgs {
		promptStr += v + " "
	}
	return strings.TrimSpace(promptStr) + ">> ", true
}
