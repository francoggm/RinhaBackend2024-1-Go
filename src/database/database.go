package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	DBClient client
)

func InitDatabase(uri string, username string, password string) error {
	conn, err := neo4j.NewDriverWithContext(uri, neo4j.NoAuth())
	if err != nil {
		return err
	}

	DBClient = &clientNeo4j{
		DB: conn,
	}

	return nil
}
