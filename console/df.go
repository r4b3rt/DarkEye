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
		if m, ok := M[ID(os.Args[1])]; ok {
			if len(os.Args) == 2 {
				m.Usage()
				return
			}
			_ = m.CompileArgs(nil, os.Args[2:])
			mContext.CmdArgs = append(mContext.CmdArgs, os.Args[1:2]...)
			m.Start(mContext.ctx)
		} else {
			usageForAll()
		}
	}
}

func interactive() {
	fmt.Printf(common.Banner)
	fmt.Println("Please use `exit` exit this program, use 'F1' 打印当前任务结果")
	p := prompt.New(
		mContext.executor,
		mContext.completer,
		prompt.OptionPrefix(">> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionLivePrefix(mContext.livePrefix),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.F1,
			Fn: func(buf *prompt.Buffer) {
				analysisRuntimeOptions.PrintCurrentTaskResult()
			}}, ),
	)

	p.Run()
}

func initializer() {
	flag.BoolVar(&mContext.Interactive, "i", false, "Launch the interactive darkEye framework")
	for _, m := range M {
		m.Init(mContext)
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

	for _, v := range Names() {
		fmt.Println(fmt.Sprintf("	%s", v))
		fmt.Println(fmt.Sprintf("		See 'df %s -h' for help", v))
	}
}
