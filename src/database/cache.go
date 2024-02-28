package database

import "time"

type ClientCache struct {
	UserID              int64
	Balance             int64
	Limit               int64
	LastTransactionUUID string
	LastTransactionDate time.Time
}

var (
	clientInfoCache         = map[int]*ClientCache{}
	clientTransactionsCache = map[int][]*Transaction{}
)

func GetClientInfoCache(id int) (*ClientCache, bool) {
	client, ok := clientInfoCache[id]
	return client, ok
}

func GetClientTransactionsCache(id int) ([]*Transaction, bool) {
	transactions, ok := clientTransactionsCache[id]
	return transactions, ok
}

func CalculateCache(id int, transactions []*Transaction) *ClientCache {
	client, ok := GetClientInfoCache(id)
	if !ok {
		client = new(ClientCache)

		client.Balance = 0
		clientInfoCache[id] = client
	}

	for _, transaction := range transactions {
		client.UserID = transaction.UserID
		client.Balance -= transaction.Value
		client.Limit = transaction.Limit
		client.LastTransactionUUID = transaction.UUID
		client.LastTransactionDate = transaction.Date
	}

	saveTransactionsCache(id, transactions)

	return client
}

func saveTransactionsCache(id int, transactions []*Transaction) {
	clientTransactions, _ := GetClientTransactionsCache(id)

	clientTransactions = append(clientTransactions, transactions...)
	if len(transactions) > 10 {
		clientTransactions = transactions[len(transactions)-10:]
	}

	clientTransactionsCache[id] = clientTransactions
}
