package database

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (c *clientNeo4j) CloseDB(ctx context.Context) error {
	return c.DB.Close(ctx)
}

func (c *clientNeo4j) VerifyConnectivity(ctx context.Context) error {
	return c.DB.VerifyConnectivity(ctx)
}

func (c *clientNeo4j) GetClientInfo(userId int64) *ClientInfo {
	query := `MATCH (t:Transaction)
						WHERE t.userId = $userId AND t.valor = 0
						
						OPTIONAL MATCH (ts:Transaction)
						WHERE ts.userId = t.userId AND ts.date > t.date
						
						RETURN t.userId as userId, t.limite as limite, sum(ts.valor) as balance, t.uuid as lastUUID, collect({uuid: ts.uuid, userId: ts.userId, date: ts.date, descricao: ts.descricao, tipo: ts.tipo, valor: ts.valor})[-10..] as transactions`

	result, _ := neo4j.ExecuteQuery(context.Background(), c.DB, query, map[string]any{"userId": userId}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

	if len(result.Records) < 1 {
		return nil
	}

	var ci *ClientInfo

	for _, info := range result.Records {
		ci = fillClientInfoRecord(info.AsMap())
	}

	return ci
}

func (c *clientNeo4j) MakeTransaction(lastSavedUUID string, userId int64, value int64, transactionType string, description string) *Transaction {
	values := map[string]any{
		"lastSavedUUID":   lastSavedUUID,
		"userId":          userId,
		"value":           value,
		"transactionType": transactionType,
		"description":     description,
	}

	query := `MATCH (t:Transaction {uuid: $lastSavedUUID})
						WHERE NOT (t)-[:NEXT]->(:Transaction)
						CREATE (t)-[:NEXT]->(nt: Transaction {
								userId: $userId, valor: $value, tipo: $transactionType, descricao: $description, date: timestamp(), uuid: randomUUID()
						})
						RETURN nt.userId as userId, nt.uuid as uuid, nt.valor as valor, nt.tipo as tipo, nt.descricao as descricao, nt.date as date`

	result, _ := neo4j.ExecuteQuery(context.Background(), c.DB, query, values, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

	if len(result.Records) < 1 {
		return nil
	}

	var transaction *Transaction

	for _, record := range result.Records {
		transaction = fillTransactionRecord(record.AsMap())
	}

	return transaction
}

func (c *clientNeo4j) IsLastTransactionUUID(uuid string) bool {
	query := `MATCH (t:Transaction {uuid: $uuid})-[:NEXT]->(n:Transaction)
						RETURN n`

	result, _ := neo4j.ExecuteQuery(context.Background(), c.DB, query, map[string]any{"uuid": uuid}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
	return len(result.Records) < 1
}

func fillClientInfoRecord(record map[string]any) *ClientInfo {
	var c ClientInfo

	id, ok := record["userId"]
	if !ok || id == nil {
		return nil
	}
	c.UserID = id.(int64)

	limit, ok := record["limite"]
	if !ok || limit == nil {
		return nil
	}
	c.Limit = limit.(int64)

	balance, ok := record["balance"]
	if !ok || balance == nil {
		return nil
	}
	c.Balance = balance.(int64)

	lastUUID, ok := record["lastUUID"]
	if !ok || lastUUID == nil {
		return nil
	}
	c.LastTransactionUUID = lastUUID.(string)

	transactions, ok := record["transactions"]
	if !ok || transactions == nil {
		return nil
	}

	c.LastTransactions = make([]*Transaction, 0)

	for _, info := range transactions.([]any) {
		transaction := fillTransactionRecord(info.(map[string]any))

		if transaction != nil {
			c.LastTransactions = append(c.LastTransactions, transaction)
			c.LastTransactionUUID = transaction.UUID
		}
	}

	return &c
}

func fillTransactionRecord(record map[string]any) *Transaction {
	var t Transaction

	id, ok := record["userId"]
	if !ok || id == nil {
		return nil
	}
	t.UserID = id.(int64)

	uuid, ok := record["uuid"]
	if !ok || uuid == nil {
		return nil
	}
	t.UUID = uuid.(string)

	value, ok := record["valor"]
	if !ok || value == nil {
		return nil
	}
	t.Value = value.(int64)

	transactionType, ok := record["tipo"]
	if !ok || transactionType == nil {
		return nil
	}
	t.TransactionType = transactionType.(string)

	description, ok := record["descricao"]
	if !ok || description == nil {
		return nil
	}
	t.Description = description.(string)

	date, ok := record["date"]
	if !ok || date == nil {
		return nil
	}
	t.Date = time.UnixMilli(date.(int64))

	return &t
}
