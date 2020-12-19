package common

import (
	"os/exec"
	"syscall"
)

func SetRLimit() {
}

func HideCmd(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
