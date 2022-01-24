//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func (c *config) requirement() error {
	if _, err := exec.LookPath("wmic.exe"); err != nil {
		return err
	}
	if _, err := exec.LookPath("schtasks.exe"); err != nil {
		return err
	}
	return nil
}

func (c *config) autorun() error {
	_, err := c.runCmd("wmic", "startup", "list", "/format:csv",
		fmt.Sprintf(">%s", filepath.Join(c.output, "autorun.txt")))
	if err != nil {
		return fmt.Errorf("autorun:%v", err.Error())
	}
	return nil
}

func (c *config) schTsk() error {
	_, err := c.runCmd("chcp", "437", "&&", "schtasks", "/query", "/fo", "csv", "/v",
		fmt.Sprintf(">%s", filepath.Join(c.output, "schTask.csv")))
	if err != nil {
		return fmt.Errorf("schTsk.runCmd:%v", err.Error())
	}

	return nil
}

func (c *config) accounts() error {
	_, err := c.runCmd("wmic", "useraccount",
		fmt.Sprintf(">%s", filepath.Join(c.output, "accounts.txt")))
	if err != nil {
		return fmt.Errorf("accounts.runCmd:%v", err.Error())
	}

	return nil
}

func (c *config) netstat() error {
	_, err := c.runCmd("netstat", "-aon",
		fmt.Sprintf(">%s", filepath.Join(c.output, "netstat.txt")))
	if err != nil {
		return fmt.Errorf("netstat.runCmd:%v", err.Error())
	}

	return nil
}

func (c *config) taskList() error {
	_, err := c.runCmd("wmic", "process", "list", "/format:csv",
		fmt.Sprintf(">%s", filepath.Join(c.output, "processes.csv")))
	if err != nil {
		return fmt.Errorf("taskList.runCmd:%v", err.Error())
	}

	return nil
}
