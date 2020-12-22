package common

import (
	"os/exec"
	"syscall"
)

//SetRLimit: add comment
func SetRLimit() {
}

//HideCmd: add comment
func HideCmd(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
