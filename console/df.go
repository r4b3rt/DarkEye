package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"runtime"
)

var (
	mContext = &RequestContext{
		CmdArgs: make([]string, 0),
	}
)

func main() {
	initializer()
	if mContext.Interactive {
		interactive()
	} else {
		if len(os.Args) == 1 {
			usageForAll()
			return
		}
		if m, ok := ModuleFuncs[moduleId(os.Args[1])]; ok {
			if len(os.Args) == 2 {
				m.usage()
				return
			}
			_ = m.compileArgs(os.Args[2:])
			m.start(mContext.ctx)
		} else {
			usageForAll()
		}
	}
}

func interactive() {
	fmt.Printf(common.Banner)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	p := prompt.New(
		mContext.executor,
		mContext.completer,
		prompt.OptionPrefix(">> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionLivePrefix(mContext.livePrefix))
	p.Run()
}

func initializer() {
	flag.BoolVar(&mContext.Interactive, "i", false, "Launch the interactive darkEye framework")
	for _, m := range ModuleFuncs {
		m.init()
	}
	flag.Usage = usageForAll
	flag.Parse()
	//初始化系统
	common.SetRLimit()
	runtime.GOMAXPROCS(runtime.NumCPU())
	mContext.ctx, mContext.cancel = context.WithCancel(context.Background())
	//  debug/pprof
	/*
		go func() {
			fmt.Println(http.ListenAndServe("localhost:10000", nil))
		}()
	*/
}

func usageForAll() {
	fmt.Println(fmt.Sprintf("Usage of %s:", os.Args[0]))
	fmt.Println("Options:")
	fmt.Println("	-i", "bool")
	fmt.Println("		Launch the darkEye framework", "(default: false)")

	for n := range ModuleFuncs {
		fmt.Println(fmt.Sprintf("	%s", n))
		fmt.Println(fmt.Sprintf("		only %s", n))
	}
}
