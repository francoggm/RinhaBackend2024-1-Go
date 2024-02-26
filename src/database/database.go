package database

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	DBContext context.Context
	DBClient  client
)

func InitDatabase(uri string, username string, password string) error {
	conn, err := neo4j.NewDriverWithContext(uri, neo4j.NoAuth())
	if err != nil {
		return err
	}

	DBContext = context.Background()

	DBClient = &clientNeo4j{
		DB: conn,
	}

	return nil
}
