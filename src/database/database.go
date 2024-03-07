package database

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type clientNeo4j struct {
	Driver             neo4j.DriverWithContext
	ExtractSession     neo4j.SessionWithContext
	TransactionSession neo4j.SessionWithContext
}

var DB *clientNeo4j

func InitDatabase(db string) error {
	driver, err := neo4j.NewDriverWithContext(fmt.Sprintf("bolt://%s:7687", db), neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		return err
	}

	ctx := context.Background()

	exs := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
		AccessMode:   neo4j.AccessModeRead,
	})

	ts := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
		AccessMode:   neo4j.AccessModeRead,
	})

	DB = &clientNeo4j{
		Driver:             driver,
		ExtractSession:     exs,
		TransactionSession: ts,
	}

	return nil
}
