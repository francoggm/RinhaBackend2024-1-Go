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
	Date            time.Time `json:"realizada_em"`
	Description     string    `json:"descricao"`
	Value           int64     `json:"valor"`
	TransactionType string    `json:"tipo"`
	Limit           int64     `json:"-"`
	UserID          int64     `json:"-"`
	UUID            string    `json:"-"`
}

type Extract struct {
	Info struct {
		Balance int64     `json:"total"`
		Date    time.Time `json:"data_extrato"`
		Limit   int64     `json:"limite"`
	} `json:"saldo"`
	Transactions []*Transaction `json:"ultimas_transacoes"`
}

type clientNeo4j struct {
	DB neo4j.DriverWithContext
}

func NewExtract(balance int64, date time.Time, limit int64, transactions []*Transaction) Extract {
	var extract Extract

	extract.Info.Balance = balance
	extract.Info.Date = date
	extract.Info.Limit = limit
	extract.Transactions = transactions

	return extract
}
