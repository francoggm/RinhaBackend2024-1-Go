package database

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/exp/slices"
)

// This implementation allows you to use multiple databases by creating a new client

type client interface {
	GetClientInfo(userId int64) *ClientInfo
	MakeTransaction(lastSavedUUID string, userId int64, value int64, transactionType string, description string) *Transaction
	IsLastTransactionUUID(uuid string) bool
	VerifyConnectivity(ctx context.Context) error
	CloseDB(ctx context.Context) error
}

type ClientInfo struct {
	UserID              int64
	Balance             int64
	Limit               int64
	LastTransactionUUID string
	LastTransactions    []*Transaction
}

type Transaction struct {
	Date            time.Time `json:"realizada_em"`
	Description     string    `json:"descricao"`
	Value           int64     `json:"valor"`
	TransactionType string    `json:"tipo"`
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
	extract.Transactions = slices.Clone(transactions)

	slices.Reverse(extract.Transactions)

	return extract
}

func (ci *ClientInfo) SetLastTransaction(transaction *Transaction) {
	ci.Balance -= transaction.Value
	ci.LastTransactionUUID = transaction.UUID

	ci.LastTransactions = append(ci.LastTransactions, transaction)
	if len(ci.LastTransactions) > 10 {
		ci.LastTransactions = ci.LastTransactions[len(ci.LastTransactions)-10:]
	}
}
