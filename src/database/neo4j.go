package database

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
)

func (c *clientNeo4j) CloseDB(ctx context.Context) error {
	return c.DB.Close(ctx)
}

func (c *clientNeo4j) VerifyConnectivity(ctx context.Context) error {
	return c.DB.VerifyConnectivity(ctx)
}

func (c *clientNeo4j) GetExtract() ([]*Extract, error) {
	return nil, nil
}

func (c *clientNeo4j) MakeTransaction(lastSavedUUID string, userId int, value int64, limit int64, transactionType string, description string) (*Transaction, error) {
	values := map[string]any{
		"lastSavedUUID":   lastSavedUUID,
		"userId":          userId,
		"value":           value,
		"limit":           limit,
		"transactionType": transactionType,
		"description":     description,
	}

	query := `MATCH (t:Transaction {uuid: $lastSavedUUID})
						WHERE NOT (t)-[:NEXT]->(:Transaction)
						CREATE (t)-[:NEXT]->(nt: Transaction {
								userId: $userId, valor: $value, limite: $limit, tipo: $transactionType, descricao: $description, date: timestamp(), uuid: randomUUID()
						})
						RETURN nt.userId as userId, nt.uuid as uuid, nt.valor as valor, nt.tipo as tipo, nt.limite as limite, nt.descricao as descricao, nt.date as date`

	result, err := neo4j.ExecuteQuery(DBContext, c.DB, query, values, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return nil, err
	}

	var transaction *Transaction

	for _, record := range result.Records {
		transaction = fillTransactionRecord(record)
	}

	return transaction, nil
}

func (c *clientNeo4j) GetAllUserTransactions(userId int) ([]*Transaction, error) {
	query := `MATCH (t:Transaction {userId: $userId}) 
	          RETURN t.userId as userId, t.uuid as uuid, t.valor as valor, t.tipo as tipo, t.limite as limite, t.descricao as descricao, t.date as date`

	result, err := neo4j.ExecuteQuery(DBContext, c.DB, query, map[string]any{"userId": userId}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return nil, err
	}

	return fillTransactions(result.Records), nil
}

func (c *clientNeo4j) GetTransactionsAfterDate(userId int64, date time.Time) ([]*Transaction, error) {
	query := `MATCH (t:Transaction)
	          WHERE t.userId = $userId AND t.date > $date
	          RETURN t.userId as userId, t.uuid as uuid, t.valor as valor, t.tipo as tipo, t.limite as limite, t.descricao as descricao, t.date as date`

	result, err := neo4j.ExecuteQuery(DBContext, c.DB, query, map[string]any{"userId": userId, "date": date.UnixMilli()}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return nil, err
	}

	return fillTransactions(result.Records), nil
}

func fillTransactionRecord(record *db.Record) *Transaction {
	m := record.AsMap()
	transaction := new(Transaction)

	id, ok := m["userId"]
	if !ok {
		return nil
	}
	transaction.UserID = id.(int64)

	uuid, ok := m["uuid"]
	if !ok {
		return nil
	}
	transaction.UUID = uuid.(string)

	value, ok := m["valor"]
	if !ok {
		return nil
	}
	transaction.Value = value.(int64)

	transactionType, ok := m["tipo"]
	if !ok {
		return nil
	}
	transaction.TransactionType = transactionType.(string)

	limit, ok := m["limite"]
	if !ok {
		return nil
	}
	transaction.Limit = limit.(int64)

	description, ok := m["descricao"]
	if !ok {
		return nil
	}
	transaction.Description = description.(string)

	date, ok := m["date"]
	if !ok {
		return nil
	}
	transaction.Date = time.UnixMilli(date.(int64))

	return transaction
}

// construct transaction struct with returned map
func fillTransactions(records []*db.Record) []*Transaction {
	var transactions = []*Transaction{}

	for _, record := range records {
		transaction := fillTransactionRecord(record)

		if transaction != nil {
			transactions = append(transactions, transaction)
		}
	}

	return transactions
}
