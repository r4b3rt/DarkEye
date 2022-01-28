package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

type threatIp struct {
	Name string
	Ip   string
	Port string
}

type pc struct {
	Name string
	Ip   string
}

type relation struct {
	Source string
	Dest   string
	How    string
}

//match (n) detach delete n
func Test_neo4j(t *testing.T) {
	threat := make(map[string]threatIp, 0)
	terminal := make(map[string]pc, 0)
	re := make(map[string]relation, 0)

	file := "/Users/b1gcat/Desktop/SL/xx/a.csv"
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c := strings.ReplaceAll(string(b), "\r", "")
	lines := strings.Split(c, "\n")
	fmt.Println("CREATE ")
	for _, line := range lines {
		if line == "" {
			break
		}
		cells := strings.Split(line, ",")
		if len(cells) != 14 {
			fmt.Println("error, not 14 vs ", len(cells), cells)
			break
		}
		//create threat
		tt := "t" + strings.ReplaceAll(cells[6], ".", "_")
		tv, ok := threat[tt]
		if !ok {
			threat[tt] = threatIp{
				Ip:   cells[6],
				Port: cells[7],
			}
		} else {
			threat[tt] = threatIp{
				Ip:   cells[6],
				Port: tv.Port + "|" + cells[7],
			}
		}
		p := "pc" + strings.ReplaceAll(cells[4], ".", "_")
		_, ok = terminal[p]
		if !ok {
			terminal[p] = pc{
				Ip: cells[4][7:],
			}
			fmt.Println(cells[4])
		}

		rv, ok := re[p+tt]
		if !ok {
			re[p+tt] = relation{
				Dest:   tt,
				Source: p,
				How:    cells[9],
			}
		} else {
			re[p+tt] = relation{
				Dest:   tt,
				Source: p,
				How:    strings.ReplaceAll(rv.How+"_"+cells[9], " ", ""),
			}
		}
	}
	//create tt
	for k, v := range threat {
		fmt.Println(fmt.Sprintf(`(%s:Threat{name:"%s",port:"%s"}),`, k, v.Ip, v.Port))
	}

	//creat pc
	for k, v := range terminal {
		fmt.Println(fmt.Sprintf(`(%s:PC{name:"%s"}),`, k, v.Ip))
	}

	//create relation
	for _, v := range re {
		fmt.Println(fmt.Sprintf(`(%s)-[:%s]->(%s),`, v.Source, v.How, v.Dest))
	}

}
