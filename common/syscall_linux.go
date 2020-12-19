package common

import "syscall"

func SetRLimit() {
	//设置max file
	rLimit := syscall.Rlimit{
		Cur: 65535,
		Max: 65535,
	}
	_ = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	_ = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

func HideCmd(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

