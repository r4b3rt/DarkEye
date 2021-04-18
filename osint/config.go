package main

import (
	"github.com/zsdevX/DarkEye/osint/graph"
)

type OsInt struct {
	Server string
	Username string
	Password string

	Graph *graph.Graph
}
