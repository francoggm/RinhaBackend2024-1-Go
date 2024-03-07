package controller

import (
	"context"
	"crebito/database"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TransactionRequest struct {
	Value       int64  `json:"valor" binding:"required"`
	Type        string `json:"tipo" binding:"required"`
	Description string `json:"descricao" binding:"required"`
}

type TransactionResponse struct {
	Balance int64 `json:"saldo"`
	Limit   int64 `json:"limite"`
}

var (
	errUserNotFound      = errors.New("user not found")
	errInsufficientLimit = errors.New("insufficient limit")
)

func MakeTransaction(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req TransactionRequest

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid fields",
		})
		return
	}

	if req.Type == "d" {
		req.Value = -1 * req.Value
	}

	c := context.Background()

	result, err := database.DB.TransactionSession.ExecuteWrite(c,
		func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(c, database.TransactionQuery,
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
				return nil, errUserNotFound
			}

			if record.AsMap()["transacao"] == nil {
				return nil, errInsufficientLimit
			}

			var res TransactionResponse

			res.Balance = record.AsMap()["saldo"].(int64)
			res.Limit = record.AsMap()["limite"].(int64)

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
