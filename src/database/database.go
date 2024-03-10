package database

import (
	"context"
	"crebito/models"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateUsers(s neo4j.SessionWithContext) {
	ctx := context.Background()

	s.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		tx.Run(ctx, CreateUsersQuery, map[string]any{})
		return nil, nil
	})
}

func GetExtract(s neo4j.SessionWithContext, id int) (any, error) {
	ctx := context.Background()

	result, err := s.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, ExtractQuery,
				map[string]any{
					"id": id,
				})
			if err != nil {
				return nil, err
			}

			record, err := result.Single(ctx)
			if err != nil {
				return nil, models.ErrUserNotFound
			}

			var res models.ExtractResponse

			res.UserInfo.Balance = record.AsMap()["saldo"].(int64)
			res.UserInfo.Limit = record.AsMap()["limite"].(int64)
			res.UserInfo.Date = time.Now()
			res.Transactions = make([]models.Transaction, 0)

			values := record.AsMap()["transacoes"].([]any)
			for _, v := range values {
				var t models.Transaction

				transaction := v.(map[string]any)

				t.Value = transaction["valor"].(int64)
				t.Type = transaction["tipo"].(string)
				t.Description = transaction["descricao"].(string)

				date := transaction["data"].(int64)
				t.Date = time.UnixMilli(date)

				res.Transactions = append(res.Transactions, t)
			}

			return res, nil
		})

	return result, err
}

func ExecuteTransaction(s neo4j.SessionWithContext, id int, req models.TransactionRequest) (any, error) {
	ctx := context.Background()

	result, err := s.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, TransactionQuery,
				map[string]any{
					"id":        id,
					"tipo":      req.Type,
					"valor":     req.Value,
					"descricao": req.Description,
				})
			if err != nil {
				return nil, err
			}

			record, err := result.Single(ctx)
			if err != nil {
				return nil, models.ErrUserNotFound
			}

			if record.AsMap()["transacao"] == nil {
				return nil, models.ErrInsufficientLimit
			}

			var res models.TransactionResponse

			res.Balance = record.AsMap()["saldo"].(int64)
			res.Limit = record.AsMap()["limite"].(int64)

			return res, nil
		})

	return result, err
}
