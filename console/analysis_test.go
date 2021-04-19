package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_analysis(t *testing.T) {
	analysisRuntimeOptions.Init(&RequestContext{})
	defer os.Remove(analysisDb)
	e1 := analysisEntity{
		Task:    "2",
		Ip:      "1.1.1.1",
		Port:    "22",
		Title:   "",
		Service: "test",
	}
	analysisRuntimeOptions.createOrUpdate(&e1)
	e2 := analysisEntity{
		Task:    "2",
		Ip:      "1.1.1.1",
		Port:    "22",
		Title:   "222",
		Service: "test",
	}
	analysisRuntimeOptions.createOrUpdate(&e2)
	e := make([]analysisEntity, 0)
	analysisRuntimeOptions.d.Raw("select * from ent").Scan(&e)
	if len(e) != 1 {
		t.Fail()
	}
	if e[0].Title != "222" {
		t.Fail()
	}

	e3 := analysisEntity{
		Task:    "1",
		Ip:      "1.1.1.1",
		Port:    "22",
		Title:   "",
		Service: "test",
	}
	analysisRuntimeOptions.createOrUpdate(&e3)
	e = make([]analysisEntity, 0)
	analysisRuntimeOptions.d.Raw("select * from ent").Scan(&e)
	fmt.Println(e)
}
