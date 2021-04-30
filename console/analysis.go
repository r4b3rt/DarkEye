package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/noborus/trdsql"
	"github.com/zsdevX/DarkEye/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"sort"
	"strings"
)

type analysisRuntime struct {
	Module
	parent *RequestContext

	d *gorm.DB

	q       string
	output  string
	flagSet *flag.FlagSet
	cmd     []string
}

var (
	analysisProgram        = "analysis"
	analysisRuntimeOptions = &analysisRuntime{
		flagSet: flag.NewFlagSet(analysisProgram, flag.ExitOnError),
	}
	analysisDb = analysisProgram + ".s3db"
)

func (a *analysisRuntime) Start(ctx context.Context) {
	d := make([]analysisEntity, 0)
	//非查询语句
	if !strings.HasPrefix(strings.ToLower(a.q), "select") {
		if err := a.d.Exec(a.q).Error; err != nil {
			common.Log("analysis.update", err.Error(), common.ALERT)
		} else {
			common.Log("analysis.update", "已经更新", common.INFO)
		}
		return
	}
	//查询语句
	ret := a.d.Raw(a.q).Scan(&d)
	if ret.Error != nil {
		return
	}
	if len(d) == 0 {
		common.Log("analysis.query", "查询无数据", common.INFO)
		return
	}
	a.OutPut(d)
}

func (a *analysisRuntime) Init(requestContext *RequestContext) {
	a.parent = requestContext
	a.flagSet.StringVar(&a.q, "sql", "select * from ent limit 1", "Sqlite3 Grammar")
	a.flagSet.StringVar(&a.output, "output-csv", "", "输出查询到文件")

	db, err := gorm.Open(sqlite.Open(analysisDb), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err.Error())
	}
	a.d = db
	err = db.AutoMigrate(&analysisEntity{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&crawler{})
	if err != nil {
		panic(err.Error())
	}
}

func (a *analysisRuntime) ValueCheck(value string) (bool, error) {
	if v, ok := analysisValueCheck[value]; ok {
		//过滤重复的命令
		if isDuplicateArg(value, a.parent.CmdArgs) {
			return false, fmt.Errorf("参数重复")
		}
		return v, nil
	}
	return false, fmt.Errorf("无此参数")
}

func (a *analysisRuntime) saveCmd(cmd []string) {

}

func (a *analysisRuntime) restoreCmd() {
	//cmd := make([]string, 0)
	a.parent.CmdArgs = a.cmd
}

func (a *analysisRuntime) CompileArgs(cmd []string) error {
	if err := a.flagSet.Parse(splitCmd(cmd)); err != nil {
		return err
	}
	a.flagSet.Parsed()
	return nil
}

func (a *analysisRuntime) Usage() {
	fmt.Println(fmt.Sprintf("Usage of %s:", analysisProgram))
	fmt.Println("Options:")
	a.flagSet.VisitAll(func(f *flag.Flag) {
		var opt = "  -" + f.Name
		fmt.Println(opt)
		fmt.Println(fmt.Sprintf("		%v (default '%v')", f.Usage, f.DefValue))
	})
}

func (a *analysisRuntime) getCrawler() ([]crawler, error) {
	c := make([]crawler, 0)
	ret := a.d.Raw(a.q).Scan(&c)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return c, nil
}

func (a *analysisRuntime) cleanCrawler() {
	if err := a.d.Exec("delete from crawler").Error; err != nil {
		common.Log("analysis.crawler.clean", err.Error(), common.ALERT)
	} else {
		common.Log("analysis.crawler.clean", "Cleaned", common.INFO)
	}
}

func (a *analysisRuntime) upInsertEnt(e *analysisEntity) {
	var n analysisEntity
	if a.d.Table("ent").Where(
		"task = ? and ip = ? and port = ? and service = ?",
		e.Task, e.Ip, e.Port, e.Service).First(&n).Error == gorm.ErrRecordNotFound {
		ret := a.d.Create(e)
		if ret.Error != nil {
			common.Log("analysis.create", ret.Error.Error(), common.ALERT)
		}
	} else {
		aJson, _ := json.Marshal(e)
		var m map[string]interface{}
		_ = json.Unmarshal(aJson, &m)
		for k, v := range m {
			switch v.(type) {
			case int32:
				if v.(int32) == 0 {
					delete(m, k)
				}
			case string:
				if v.(string) == "" {
					delete(m, k)
				}
			default:
				delete(m, k)
			}
		}
		a.d.Model(e).Where(
			"task = ? and ip = ? and port = ? and service = ?",
			e.Task, e.Ip, e.Port, e.Service).Updates(
			m)
	}
}

func (a *analysisRuntime) upInsertCrawler(c *crawler) {
	var n analysisEntity
	if a.d.Table("crawler").Where(
		"target = ? and url = ? and method = ?",
		c.Target, c.Url, c.Method).First(&n).Error == gorm.ErrRecordNotFound {
		ret := a.d.Create(c)
		if ret.Error != nil {
			common.Log("analysis.create", ret.Error.Error(), common.ALERT)
		}
	} else {
		if c.Data != "" {
			a.d.Model(c).Where(
				"target = ? and url = ? and method = ?",
				c.Target, c.Url, c.Method).Updates(
				c)
		}
	}
}

func (a *analysisRuntime) Var(condition string, field string) ([]string, error) {
	if !strings.HasPrefix(field, "$") {
		return nil, fmt.Errorf(field, "非变量")
	}
	field = strings.ToLower(strings.TrimPrefix(field, "$"))
	e := make([]analysisEntity, 0)
	sql := fmt.Sprintf("select %s from ent", field)
	a.d.Raw(sql).Scan(&e)
	ret := make([]string, 0)
	for _, v := range e {
		aJson, _ := json.Marshal(v)
		var m map[string]interface{}
		_ = json.Unmarshal(aJson, &m)
		if value, ok := m[field]; ok {
			ret = append(ret, fmt.Sprint(value))
		}
	}
	return ret, nil
}

func (a *analysisRuntime) PrintCurrentTaskResult() {
	e := make([]analysisEntity, 0)
	sql := fmt.Sprintf(
		"select task,ip,port,service,country,isp,weak_account,title,url,netbios from ent where task=%d",
		a.parent.taskId)
	a.d.Raw(sql).Scan(&e)
	a.OutPut(e)
}

func (a *analysisRuntime) OutPut(d []analysisEntity) {
	sort.Slice(d, func(i, j int) bool {
		return d[j].Ip > d[i].Ip
	})

	//查询结果写入内存
	jsonString, _ := json.Marshal(d)
	r := bytes.NewBuffer(jsonString)
	importer, err := trdsql.NewBufferImporter("any", r, trdsql.InFormat(trdsql.JSON))
	if err != nil {
		common.Log(a.parent.CmdArgs[0], err.Error(), common.FAULT)
		return
	}
	//查询输出table格式
	writer := trdsql.NewWriter(trdsql.OutFormat(trdsql.AT))
	trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
	trd.Driver = "sqlite3"
	err = trd.Exec("select * from any")
	if err != nil {
		common.Log(a.parent.CmdArgs[0], err.Error(), common.FAULT)
		return
	}
	//查询导出文件
	if a.output != "" {
		fp, err := os.Create(a.output)
		if err != nil {
			common.Log(a.parent.CmdArgs[0]+".output", err.Error(), common.FAULT)
		}
		defer fp.Close()
		writer := trdsql.NewWriter(trdsql.OutFormat(trdsql.CSV),
			trdsql.OutHeader(true),
			trdsql.OutStream(fp))
		trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
		err = trd.Exec("select * from any")
		if err != nil {
			common.Log(a.parent.CmdArgs[0], err.Error(), common.FAULT)
			return
		}
	}
}
