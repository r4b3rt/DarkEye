package main

import "syscall"

func setRLimit() {
	//设置max file
	rLimit := syscall.Rlimit{
		Cur: 65535,
		Max: 65535,
	}
	_ = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	_ = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}
