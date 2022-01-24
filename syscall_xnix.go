//go:build !windows

package main

import (
	"github.com/sirupsen/logrus"
	"syscall"
)

func setNoFiles() {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logrus.Error("setNoFiles.get:", err.Error())
		return
	}
	logrus.Info("current max open files:", rLimit)

	rLimit.Max = 65535
	rLimit.Cur = 65535

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logrus.Error("setNoFiles.set:", err.Error())
		return
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logrus.Error("setNoFiles.get:", err.Error())
		return
	}
	logrus.Info("reset max open files:", rLimit)
}
