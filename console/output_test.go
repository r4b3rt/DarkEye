package main

import "testing"

func TestSuperScanRuntime_OutPut(t *testing.T) {
	initializer()

	superScanRuntimeOptions.result = []analysisEntity{
		{
			Ip: "1.1.1.1",
		},
		{
			Ip: "1.1.1.2",
		},
	}
	s := *superScanRuntimeOptions
	s.OutPut()
}
