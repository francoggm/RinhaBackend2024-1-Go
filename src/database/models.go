package database

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// This implementation allows you to use multiple databases by creating a new client

type client interface {
	GetExtract() ([]*Extract, error)
	MakeTransaction(lastSavedUUID string, userId int, value int64, limit int64, transactionType string, description string) (*Transaction, error)
	GetAllUserTransactions(userId int) ([]*Transaction, error)
	GetTransactionsAfterDate(userId int64, date time.Time) ([]*Transaction, error)
	VerifyConnectivity(ctx context.Context) error
	CloseDB(ctx context.Context) error
}

type Transaction struct {
	Date            time.Time
	Description     string
	Value           int64
	Limit           int64
	TransactionType string
	UserID          int64
	UUID            string
}

type Extract struct {
}

type clientNeo4j struct {
	DB neo4j.DriverWithContext
}
