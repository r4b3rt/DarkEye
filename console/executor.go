package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

func (ctx *RequestContext) executor(in string) {
	in = strings.TrimPrefix(strings.TrimSpace(in), "?")
	blocks := strings.Split(in, " ")
	switch strings.ToLower(blocks[0]) {
	case "":
		break
	case "stop":
		ctx.stop()
	case "exit":
		ctx.stop()
		fmt.Println("Bye Bye!")
		os.Exit(0)
	case "cd":
		if ctx.checkRun() {
			return
		}
		if len(blocks) < 2 {
			ctx.CmdArgs = make([]string, 0)
			return
		}
		if blocks[1] == ".." {
			ctx.returnCmd()
			return
		}
		//处理主命令
		if cmdValid(blocks[1]) {
			ctx.CmdArgs = append(ctx.CmdArgs, blocks[1])
		} else {
			fmt.Println(fmt.Sprintf("'%s' Not support Command", blocks[1]))
		}
		return
	default:
		if ctx.checkRun() {
			return
		}
		switch len(ctx.CmdArgs) == 0 {
		case true:
			if cmdValid(blocks[0]) {
				ctx.CmdArgs = append(ctx.CmdArgs, blocks[0])
				return
			}
		case false:
			ctx.runCmd(blocks)
			return
		}
		fmt.Println(fmt.Sprintf("You Should try input: 'cd %s'", blocks[0]))
	}
}

func (ctx *RequestContext) runCmd(args []string) {
	switch strings.ToLower(args[0]) {
	case "exploit":
		if err := ModuleFuncs[moduleId(ctx.CmdArgs[0])].compileArgs(ctx.CmdArgs[1:]); err != nil {
			fmt.Println("Err:", err.Error())
			return
		}
		ctx.running.Store(true)
		go func() {
			color.Green("%s %s", ctx.CmdArgs[0], "Running!")
			ModuleFuncs[moduleId(ctx.CmdArgs[0])].start(mContext.ctx)
			ctx.running.Store(false)
			color.Green("%s %s", ctx.CmdArgs[0], "Done!")
		}()
	default:
		if len(args) == 1 { //此时检查是否需要参数
			if noVar, ok := ModuleFuncs[moduleId(ctx.CmdArgs[0])].valueCheck[args[0]]; !ok {
				fmt.Println("Err: not support", "'"+args[0]+"'")
				return
			} else {
				if !noVar {
					fmt.Println("Err:", "'"+args[0]+"'", "need value")
					return
				}
			}
		}
		cmd := ""
		for _, v := range args {
			cmd += v + " "
		}
		ctx.CmdArgs = append(ctx.CmdArgs, strings.TrimSpace(cmd))
	}
}

func runShell(cmd *exec.Cmd) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdoutPipe.Close()

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() { // 命令在执行的过程中, 实时地获取其输出
			fmt.Println(scanner.Bytes())
		}
	}()

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func splitCmd(cmd []string) []string {
	ret := make([]string, 0)
	for _, v := range cmd {
		s := strings.SplitN(v, " ", 2)
		ret = append(ret, s...)
	}
	return ret
}

func (ctx *RequestContext) returnCmd() {
	ctx.CmdArgs = ctx.CmdArgs[:len(ctx.CmdArgs)-1]
	return
}

func (ctx *RequestContext) stop() {
	if !ctx.running.Load() {
		fmt.Println("Stopped")
		return
	}
	ctx.cancel()
	ctx.ctx, ctx.cancel = context.WithCancel(context.Background())
	fmt.Println("Waiting to stop")
}

func (ctx *RequestContext) checkRun() bool {
	if ctx.running.Load() {
		fmt.Println(ctx.CmdArgs[0], "is running")
		return true
	}
	return false
}

func cmdValid(cmd string) bool {
	for _, v := range mSuggestions {
		if strings.ToLower(v.Text) == strings.ToLower(cmd) {
			return true
		}
	}
	return false
}
