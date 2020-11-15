package poc

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"path/filepath"
	"strings"
)

func (p *Poc) Check() {
	defer func() {
		p.ErrChannel <- common.LogBuild("poc",
			fmt.Sprintf("POC检查结束"), common.INFO)
	}()
	fi, err := os.Stat(p.FileName)
	if err != nil {
		p.ErrChannel <- common.LogBuild("poc",
			fmt.Sprintf("读取目录或文件错误:%s", err.Error()), common.FAULT)
		return
	}
	if !fi.IsDir() {
		p.CheckByFileName(p.FileName)
		return
	} else {
		//打开目录
		f, err := os.OpenFile(p.FileName, os.O_RDONLY, os.ModeDir)
		if err != nil {
			p.ErrChannel <- common.LogBuild("poc",
				fmt.Sprintf("打开目录错误:%s", err.Error()), common.FAULT)
			return
		}
		defer f.Close()

		//读取目录
		rd, err := f.Readdir(-1)
		if err != nil {
			p.ErrChannel <- common.LogBuild("poc",
				fmt.Sprintf("读取目录错误:%s", err.Error()), common.FAULT)
			return
		}
		for _, rdi := range rd {
			if rdi.IsDir() {
				continue
			}
			if common.ShouldStop(&p.Stop) {
				break
			}
			if strings.HasSuffix(rdi.Name(), "yml") {
				p.CheckByFileName(filepath.Join(p.FileName, rdi.Name()))
			}
		}
	}
}

func (p *Poc) CheckByFileName(pocName string) {
	urls := strings.Split(p.Urls, ",")
	for _, myUrl := range urls {
		if common.ShouldStop(&p.Stop) {
			break
		}
		myUrl = strings.TrimSpace(strings.TrimRight(myUrl, "/"))
		result, err := p.doCheck(pocName, myUrl)
		if err != nil {
			p.ErrChannel <- common.LogBuild("poc",
				fmt.Sprintf("CheckByFile:/%s/%s", pocName, err.Error()), common.FAULT)
		}
		infoLevel := common.INFO
		if result {
			infoLevel = common.ALERT
		}
		p.ErrChannel <- common.LogBuild("poc",
			fmt.Sprintf("%s %s:%v", myUrl, filepath.Base(pocName), result), infoLevel)
	}
	return
}
