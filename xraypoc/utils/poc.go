package main

import (
	"flag"
	"fmt"
	"github.com/zsdevX/DarkEye/xraypoc"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type poc struct {
	Name string
	Data string
}

var (
	dictionaryPoc    = "../db_poc"
	pocs             = make([]poc, 0)
	testUrl          = flag.String("test-url", "http://127.0.0.1", "测试url")
	test             = flag.Bool("test", false, "测试poc")
	testPoc          = flag.String("test-poc", "poc.yml", "测试的poc")
	mPocReverse      = flag.String("reverse-url", "qvn0kc.ceye.io", "CEye 标识")
	mPocReverseCheck = flag.String("reverse-check-url", "http://api.ceye.io/v1/records?token=066f3d242991929c823ac85bb60f4313&type=http&filter=", "CEye API")
)

func main() {
	flag.Parse()
	if *test {
		pocTest()
		return
	}
	//打开目录
	f, err := os.OpenFile(dictionaryPoc, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return
	}
	defer f.Close()

	//读取目录
	rd, err := f.Readdir(-1)
	if err != nil {
		return
	}
	for _, rdi := range rd {
		if rdi.IsDir() {
			continue
		}
		if strings.HasSuffix(rdi.Name(), "yml") {
			d, _ := ioutil.ReadFile(filepath.Join(dictionaryPoc, rdi.Name()))
			pocs = append(pocs, poc{
				Name: rdi.Name(),
				Data: string(d),
			})
		}
	}
	genPocCode()
}

func genPocCode() {
	output, _ := os.OpenFile(filepath.Join(dictionaryPoc, "poc.go"),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	defer output.Close()

	_, _ = output.WriteString("package db_poc\n")
	_, _ = output.WriteString(`type Poc struct {
	Name string
	Data string
}`)
	_, _ = output.WriteString("\n\n")
	_, _ = output.WriteString("var POCS =[]Poc{\n")
	for _, p := range pocs {
		_, _ = output.WriteString("{\n")
		_, _ = output.WriteString(fmt.Sprintf(`Name:"%s",`, p.Name))
		_, _ = output.WriteString("\n")
		_, _ = output.WriteString(fmt.Sprintf("Data:`%s`,", p.Data))
		_, _ = output.WriteString("},\n")
	}
	_, _ = output.WriteString("}\n")
}

func pocTest() {
	xAry := xraypoc.XArYPoc{
		ReverseUrlCheck: *mPocReverse,
		ReverseUrl:      *mPocReverseCheck,
	}
	ok, err := xAry.Check(nil, *testPoc, *testUrl)
	fmt.Println(fmt.Sprintf("%v:%v", ok, err))
}
