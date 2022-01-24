package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func (c *config) gatherInfo() {
	if c.noGather {
		c.logger.Warn("gather disabled")
		return
	}
	c.logger.Info("start gathering ...")
	defer c.logger.Info("gather done!")

	if err := c.requirement(); err != nil {
		c.logger.Error("gatherInfo.requirement:", err.Error())
		return
	}

	err := c.autorun()
	if err != nil {
		c.logger.Error("gatherInfo.:", err.Error())
		return
	}

	err = c.schTsk()
	if err != nil {
		c.logger.Error("gatherInfo.:", err.Error())
		return
	}

	err = c.accounts()
	if err != nil {
		c.logger.Error("gatherInfo.:", err.Error())
		return
	}

	err = c.taskList()
	if err != nil {
		c.logger.Error("gatherInfo.:", err.Error())
		return
	}

	err = c.netstat()
	if err != nil {
		c.logger.Error("gatherInfo.:", err.Error())
		return
	}
}

func (*config) runCmd(cmd ...string) ([]byte, error) {
	var c *exec.Cmd
	switch runtime.GOOS {
	case "linux":
	case "darwin":
		c = exec.Command("sh", "-c", strings.Join(cmd, " "))
	case "windows":
		c = exec.Command("CMD", "/C", strings.Join(cmd, " "))
	default:
		return nil, fmt.Errorf("runcmd not support arch")
	}
	output, err := c.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("err=%v", err.Error())
	}
	return output, nil
}
