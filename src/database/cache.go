package database

var clientInfoCache = make(map[int64]*ClientInfo)

func GetClientInfoCache(id int64) (*ClientInfo, bool) {
	client, ok := clientInfoCache[id]
	return client, ok
}

func SetClientInfoCache(info *ClientInfo) {
	clientInfoCache[info.UserID] = info
}

func CalculateCache(id int64, transaction *Transaction) {
	client, _ := GetClientInfoCache(id)
	client.SetLastTransaction(transaction)
}
