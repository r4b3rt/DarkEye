package graph

import "github.com/neo4j/neo4j-go-driver/neo4j"


const (
	AddNode     = 1
	AddRelation = 2
	CleanAll    = 3
)

type Graph struct {
	Driver   neo4j.Driver
	Server   string
	Username string
	Password string
}

func New(dbUrl, username, password string) (*Graph, error) {
	d, err := neo4j.NewDriver(dbUrl,
		neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	neo := &Graph{
		Driver:   d,
		Username: username,
		Password: password,
		Server:   dbUrl,
	}
	return neo, nil
}
