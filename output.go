package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func (c *config) output(l ...interface{}) {
	if !c.bar {
		fmt.Println(l...)
	}
	writeTo(c.outfile, l)
}

func writeTo(filename string, a ...interface{}) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error("writeTo:", err.Error())
		return
	}
	defer f.Close()
	fmt.Fprintln(f, a)
}
