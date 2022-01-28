package main

import (
	"fmt"
	"github.com/b1gcat/go-libraries/utils"
	"strconv"
	"strings"
	"testing"
)

func Test_match(t *testing.T) {
	a, _ := utils.LoadListFromFile("dist/a.txt")
	fmt.Println(len(a))
	i := 0
	for i <= 5 {
		i++
		b, _ := utils.LoadListFromFile("dist/" + strconv.Itoa(i) + ".csv")
		for _, a0 := range a {
			for _, b0 := range b {
				if strings.Contains(b0, a0) {
					fmt.Println(i, a0, b0)
				}
			}
		}

	}
}
