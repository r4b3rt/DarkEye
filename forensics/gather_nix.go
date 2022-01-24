//go:build !windows

package main

import "fmt"

func (c *config) requirement() error {
	return fmt.Errorf("not support")
}

func (c *config) autorun() error {
	return fmt.Errorf("not support")
}

func (c *config) schTsk() error {
	return fmt.Errorf("not support")
}

func (c *config) accounts() error {
	return fmt.Errorf("not support")
}

func (c *config) netstat() error {
	return fmt.Errorf("not support")
}

func (c *config) taskList() error {
	return fmt.Errorf("not support")
}
