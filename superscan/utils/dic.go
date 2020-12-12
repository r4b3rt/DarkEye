//字典集成工具
//将dic目录中的字典*.txt转为*.go

package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var (
	dictionary = "../dic"
)

func main() () {
	f, _ := os.OpenFile(dictionary, os.O_RDONLY, os.ModeDir)
	defer f.Close()

	//读取目录
	rd, _ := f.Readdir(-1)
	for _, rdi := range rd {
		if rdi.IsDir() {
			continue
		}
		if !strings.HasSuffix(rdi.Name(), "txt") {
			continue
		}
		genCode(rdi.Name())
	}
}

func genCode(filename string) {
	file, _ := os.Open(filepath.Join(dictionary, filename))
	tag := strings.TrimSuffix(filename, ".txt")
	output, _ := os.OpenFile(filepath.Join(dictionary, tag+".go"),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	defer output.Close()
	defer file.Close()
	output.WriteString("package dic\n")
	output.WriteString("var " + strings.ToUpper(tag) + "= []string{\n")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		one := scanner.Text()
		if strings.HasPrefix(one, "#") {
			continue
		}
		one = strings.TrimSpace(one)
		one = strings.Trim(one, "\r\n")
		if one == "" {
			continue
		}
		output.WriteString(`"` + one + `"` + ",\n")
	}
	output.WriteString("}")
}
