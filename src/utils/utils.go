package utils

func CanMakeTransaction(transactionType string, value int64, balance int64, limit int64) bool {
	return (transactionType == "c") || ((transactionType == "d") && (balance-value > -limit))
}
