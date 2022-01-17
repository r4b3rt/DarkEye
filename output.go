package main

import "fmt"

func (c *config) output(l ...interface{}) {
	fmt.Println(l...)
}
