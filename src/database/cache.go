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
	clientsCache = map[int]*ClientCache{}
)

func GetClientCache(id int) (*ClientCache, bool) {
	client, ok := clientsCache[id]
	return client, ok
}

func CalculateCache(id int, transactions []*Transaction) *ClientCache {
	client, ok := GetClientCache(id)
	if !ok {
		client = new(ClientCache)

		client.Balance = 0
		clientsCache[id] = client
	}

	for _, transaction := range transactions {
		client.UserID = transaction.UserID
		client.Balance -= transaction.Value
		client.Limit = transaction.Limit
		client.LastTransactionUUID = transaction.UUID
		client.LastTransactionDate = transaction.Date
	}

	return client
}
