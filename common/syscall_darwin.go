package common

import (
	"os/exec"
	"syscall"
)

//SetRLimit add comment
func SetRLimit() {
	//设置max file
	rLimit := syscall.Rlimit{
		Cur: 65535,
		Max: 65535,
	}
	_ = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	_ = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

//HideCmd add comment
func HideCmd(c *exec.Cmd) {

}
