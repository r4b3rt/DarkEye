package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	runShellOutput = "shell.out"
)

func (ctx *RequestContext) executor(in string) {
	in = strings.TrimPrefix(strings.TrimSpace(in), "?")
	blocks := strings.Split(in, " ")
	switch strings.ToLower(blocks[0]) {
	case "stop":
		ctx.stop()
	case "exit":
		ctx.stop()
		common.Log("executor", "Bye Bye!", common.INFO)
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
		return
	case "":
		break
	default:
		if ctx.checkRun() {
			return
		}
		switch len(ctx.CmdArgs) == 0 {
		case true:
			if cmdValid(blocks[0]) {
				ctx.CmdArgs = append(ctx.CmdArgs, blocks[0])
				//恢复历史命令
				ctx.CmdArgs = append(ctx.CmdArgs, M[ID(ctx.CmdArgs[0])].restoreCmd()...)
				return
			} else {
				common.Log("executor", fmt.Sprintf("'%s' Not support Command", blocks[1]), common.FAULT)
			}
		case false:
			ctx.runCmd(blocks)
			return
		}
		common.Log("executor", fmt.Sprintf("You Should try input: 'cd %s'", blocks[0]), common.ALERT)
	}
}

func (ctx *RequestContext) runCmd(args []string) {
	switch strings.ToLower(args[0]) {
	case "exploit":
		if err := M[ID(ctx.CmdArgs[0])].CompileArgs(ctx.CmdArgs[1:], nil); err != nil {
			common.Log("executor.exploit", err.Error(), common.FAULT)
			return
		}
		ctx.running.Store(true)
		go func() {
			common.Log(ctx.CmdArgs[0], "Running!", common.INFO)
			M[ID(ctx.CmdArgs[0])].Start(mContext.ctx)
			ctx.running.Store(false)
			common.Log(ctx.CmdArgs[0], "Done!", common.INFO)
		}()
	default:
		//此时检查参数
		if noVar, err := M[ID(ctx.CmdArgs[0])].ValueCheck(args[0]); err != nil {
			common.Log("executor.runCmd", "'"+args[0]+"' "+err.Error(), common.FAULT)
			return
		} else {
			if !noVar && len(args) == 1 {
				common.Log("executor.runCmd", "Err:"+"'"+args[0]+"'"+" need value", common.FAULT)
				return
			}
		}
		cmd := ""
		for _, v := range args {
			cmd += v + " "
		}
		ctx.CmdArgs = append(ctx.CmdArgs, strings.TrimSpace(cmd))
		if len(ctx.CmdArgs) > 1 {
			//保存
			M[ID(ctx.CmdArgs[0])].saveCmd(ctx.CmdArgs)
		}
	}
}

func runShell(name string, cmd *exec.Cmd, saveOutput bool) error {
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() {
		_ = stdoutPipe.Close()
		_ = stderrPipe.Close()
	}()
	common.Log("runShell."+name, "start", common.INFO)
	if err := cmd.Start(); err != nil {
		return err
	}
	if saveOutput {
		f, err := os.OpenFile(runShellOutput, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		go func() {
			_, _ = io.Copy(f, stdoutPipe)
			_ = f.Close()
		}()
	} else {
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() { // 命令在执行的过程中, 实时地获取其输出
				fmt.Println(scanner.Text())
			}
		}()
	}
	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() { // 命令在执行的过程中, 实时地获取其输出
		fmt.Println(scanner.Text())
	}
	return cmd.Wait()
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
	if len(ctx.CmdArgs) > 1 {
		M[ID(ctx.CmdArgs[0])].saveCmd(ctx.CmdArgs)
	}
	return
}

func (ctx *RequestContext) stop() {
	if !ctx.running.Load() {
		common.Log("executor.runCmd", "Stopped", common.INFO)
		return
	}
	ctx.cancel()
	ctx.ctx, ctx.cancel = context.WithCancel(context.Background())
	common.Log("executor.runCmd", "Waiting to stop", common.INFO)
}

func (ctx *RequestContext) checkRun() bool {
	if ctx.running.Load() {
		common.Log("executor.checkRun", ctx.CmdArgs[0]+" is running", common.INFO)
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

func cmdSave(cmd []string) []string {
	if len(cmd) == 1 {
		return []string{}
	}
	cmd0 := make([]string, len(cmd)-1)
	copy(cmd0, cmd[1:])
	return cmd0
}

func cmdRestore(a []string) []string {
	if len(a) == 0 {
		return []string{}
	}
	cmd := make([]string, len(a))
	copy(cmd, a)
	return cmd
}
