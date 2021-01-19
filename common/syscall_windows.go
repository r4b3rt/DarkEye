package common

import (
	"os/exec"
	"syscall"
)

//SetRLimit add comment
func SetRLimit() {
}

//HideCmd add comment
func HideCmd(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

//AutoRun add comment
func AutoRun(app string, run bool) error {
	f := func(k registry.Key, run bool) error {
		if run {
			if err := k.SetStringValue(app, SysExecPath(app+".exe")); err != nil {
				return err
			}
		} else {
			k.DeleteValue(app)
		}
		return nil
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if err = f(k, run); err != nil {
		return err
	}
	//FIXME: win10
	k1, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return nil
	}
	defer k1.Close()
	f(k1, run)
	return nil
}
