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
	sort.Slice(d, func(i, j int) bool {
		return d[j].Ip > d[i].Ip
	})

	jsonString, _ := json.Marshal(d)
	r := bytes.NewBuffer(jsonString)
	importer, err := trdsql.NewBufferImporter("any", r, trdsql.InFormat(trdsql.JSON))
	if err != nil {
		common.Log(a.parent.CmdArgs[0], err.Error(), common.FAULT)
		return
	}

	writer := trdsql.NewWriter(trdsql.OutFormat(trdsql.AT))
	trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
	trd.Driver = "sqlite3"
	err = trd.Exec("select * from any")
	if err != nil {
		common.Log(a.parent.CmdArgs[0], err.Error(), common.FAULT)
		return
	}

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

func (a *analysisRuntime) createOrUpdate(e *analysisEntity) {
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
