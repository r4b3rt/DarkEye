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
	flag.StringVar(&osIntRuntimeOptions.Server, "graph-server", "bolt://localhost:7687", "server")
	flag.StringVar(&osIntRuntimeOptions.Username, "graph-username", "neo4j", "username")
	flag.StringVar(&osIntRuntimeOptions.Password, "graph-password", "changeit", "password")
	flag.Parse()
	var err error
	osIntRuntimeOptions.Graph, err = graph.New(osIntRuntimeOptions.Server, osIntRuntimeOptions.Username, osIntRuntimeOptions.Password)
	if err != nil {
		panic(err)
	}
}
