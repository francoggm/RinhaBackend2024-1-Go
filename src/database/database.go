package database

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var DBClient client

func InitDatabase(db string) error {
	conn, err := neo4j.NewDriverWithContext(fmt.Sprintf("bolt://%s:7687", db), neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		return err
	}

	DBClient = &clientNeo4j{
		DB: conn,
	}

	return nil
}
