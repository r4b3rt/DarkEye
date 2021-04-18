package graph

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func (n *Graph) AdjustItem(operation int, D map[string]interface{}) (interface{}, error) {
	// Sessions are short-lived, cheap to create and NOT thread safe. Typically create one or more sessions
	// per request in your web application. Make sure to call Close on the session when done.
	// For multi-database support, set sessionConfig.DatabaseName to requested database
	// Session config will default to write mode, if only reads are to be used configure session for
	// read mode.
	session := n.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	switch operation {
	case AddNode:
		return session.Run("MERGE (n:Person { id_str: $id, name: $name, social: $social })", D)
	case AddRelation:
		cyPer := `MATCH (n1:Person),(n2:Person) 
						WHERE ` + D["condition"].(string) + ` MERGE (n1)-[r:Follow]->(n2)`
		return session.Run(cyPer, D)
	case CleanAll:
		return session.Run("match (n) detach delete n", nil)
	default:
		return nil, fmt.Errorf("Not implementation yet! ")
	}
}
