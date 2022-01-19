package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type resultScan struct {
	Output []interface{}
}

func (c *config) output(l ...interface{}) {
	r := resultScan{
		Output: l,
	}
	body, err := json.Marshal(&r)
	if err != nil {
		logrus.Error("output.Marshal:", err.Error())
		return
	}
	if !c.bar {
		fmt.Println(string(body))
	}
	writeTo(c.outfile, string(body))
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
