package main

import (
	"context"
	"github.com/c-bata/go-prompt"
	"strings"
)

type Module interface {
	Start(ctx context.Context)
	Init(*RequestContext)
	ValueCheck(string) (bool, error)
	CompileArgs([]string, []string) error
	Completer([]string) []prompt.Suggest
	saveCmd(cmd []string)
	restoreCmd() []string
	Usage()
}

var (
	M = map[string]Module{
		ID(zoomEye):         zoomEyeRuntimeOptions,
		ID(xRayProgram):     xRayRuntimeOptions,
		ID(analysisProgram): analysisRuntimeOptions,
		ID(superScan):       superScanRuntimeOptions,
	}
)

func New(name string) Module {
	if v, ok := M[name]; ok {
		return v
	}
	return nil
}

func ID(m string) string {
	return strings.ToLower(m)
}

func Names() []string {
	var names []string
	for k := range M {
		names = append(names, k)
	}
	return names
}
