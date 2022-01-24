package main

import (
	"bytes"
	"strings"
)

func (c *config) parseLines(input []byte) []string {
	input = bytes.ReplaceAll(input, []byte("\r"), []byte(""))
	r := strings.Split(string(input), "\n")
	return r
}
