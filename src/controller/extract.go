package controller

import (
	"context"
	"crebito/database"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Transaction struct {
	Value       int64     `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Date        time.Time `json:"realizada_em"`
}

type Info struct {
	Balance int64     `json:"total"`
	Limit   int64     `json:"limite"`
	Date    time.Time `json:"data_extrato"`
}

type ExtractResponse struct {
	UserInfo     Info          `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
}

func GetExtract(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	c := context.Background()

	result, err := database.DB.ExtractSession.ExecuteRead(c,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(c, database.ExtractQuery,
				map[string]any{
					"id": id,
				})
			if err != nil {
				return nil, err
			}

			record, err := result.Single(c)
			if err != nil {
				return nil, errUserNotFound
			}

			var res ExtractResponse

			res.UserInfo.Balance = record.AsMap()["saldo"].(int64)
			res.UserInfo.Limit = record.AsMap()["limite"].(int64)
			res.UserInfo.Date = time.Now()

			values := record.AsMap()["transacoes"].([]any)
			for _, v := range values {
				var t Transaction

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

	if err != nil {
		switch {
		case errors.Is(err, errUserNotFound):
			ctx.Status(http.StatusNotFound)
		default:
			ctx.Status(http.StatusUnprocessableEntity)
		}

		return
	}

	ctx.JSON(http.StatusOK, result)
}
