package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/b1gcat/DarkEye/common"
	"github.com/noborus/trdsql"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	programName      = "DarkEye"
	defaultOutputDir = "." //配置、输出目录

	webHook      = "" //第三方webhook地址
	localWebHook = "" //本地webhook地址 local forward => webhook
)

func main() {
	runApp()
}

func output(d []byte, f string) error {
	r := bytes.NewBuffer(d)
	importer, err := trdsql.NewBufferImporter("any", r, trdsql.InFormat(trdsql.JSON))
	if err != nil {
		return err
	}
	fp, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fp.Close()

	writer := trdsql.NewWriter(trdsql.OutFormat(trdsql.CSV),
		trdsql.OutHeader(true),
		trdsql.OutStream(fp))
	trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
	err = trd.Exec("select * from any")
	if err != nil {
		return err
	}
	return nil
}

func localWebHookInit() {
	if localWebHook != "" {
		//本地服务已经启动
		return
	}
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		widgets.QMessageBox_Information(nil, "创建webhook错误", "创建webhook受限，无法监听本地端口",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	localWebHook = "http://127.0.0.1:" + strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	go func(l net.Listener) {
		http.HandleFunc("/", forwardWebHook)
		s := http.Server{}
		s.Serve(l)
		defer l.Close()
	}(l)
}

func forwardWebHook(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		widgets.QMessageBox_Information(nil, "错误", "读取localWebHook内容失败:"+err.Error(),
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}

	//格式化
	if !bytes.Contains(body, []byte(`"type":"web_vuln"`)) {
		//不支持的格式
		return
	}
	var out bytes.Buffer
	err = json.Indent(&out, body, "", "\t")
	if err != nil {
		return
	}
	myBody := "```\n" + out.String() + "\n```"
	//构造第三方格式豹纹
	//https://sctapi.ftqq.com/****************.send?title=messagetitle&desp=messagecontent
	if strings.Contains(webHook, "sctapi.ftqq.com") {
		myBody = fmt.Sprintf("title=darkEye&desp=%s", myBody)
	}
	m := common.HttpRequest{
		Method:  "POST",
		Url:     webHook,
		TimeOut: time.Duration(10),
		Body:    []byte(myBody),
		Headers: map[string]string{
			"Content-type": "application/x-www-form-urlencoded",
		},
	}
	m.Go()
}
