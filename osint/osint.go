package main

import (
	"flag"
	"github.com/zsdevX/DarkEye/osint/graph"
)

var (
	osIntRuntimeOptions = OsInt{}
)

func main() {
	initialize()
	restFul()
}

func initialize() {
	flag.StringVar(&osIntRuntimeOptions.Server, "graph-server", "bolt://192.168.1.46:7687", "server")
	flag.StringVar(&osIntRuntimeOptions.Username, "graph-username", "neo4j", "username")
	flag.StringVar(&osIntRuntimeOptions.Password, "graph-password", "88888888", "password")
	flag.Parse()
	var err error
	osIntRuntimeOptions.Graph, err = graph.New(osIntRuntimeOptions.Server, osIntRuntimeOptions.Username, osIntRuntimeOptions.Password)
	if err != nil {
		panic(err)
	}
}
