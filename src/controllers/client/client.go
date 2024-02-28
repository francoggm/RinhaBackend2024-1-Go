package client

import (
	"crebito/database"
	"crebito/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionBody struct {
	Value       int64  `json:"valor" binding:"required"`
	Type        string `json:"tipo" binding:"required"`
	Description string `json:"descricao" binding:"required"`
}

func CreateRoutes(c *gin.RouterGroup) {
	c.GET("/:id/extrato", getExtract)
	c.POST("/:id/transacoes", makeTransaction)
}

func getExtract(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	client, ok := database.GetClientInfoCache(id)
	if !ok {
		// user not found
		if !database.DBClient.FindUser(id) {
			ctx.Status(http.StatusNotFound)
			return
		}

		// cache doesnt exists, try to get user from db
		transactions, err := database.DBClient.GetAllUserTransactions(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		client := database.CalculateCache(id, transactions)

		transactions, _ = database.GetClientTransactionsCache(id)
		extract := database.NewExtract(client.Balance, time.Now(), client.Limit, transactions)

		ctx.JSON(http.StatusOK, extract)
		return
	}

	transactions, err := database.DBClient.GetTransactionsAfterDate(client.UserID, client.LastTransactionDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(transactions) > 0 {
		client = database.CalculateCache(id, transactions)
	}

	transactions, _ = database.GetClientTransactionsCache(id)
	extract := database.NewExtract(client.Balance, time.Now(), client.Limit, transactions)

	ctx.JSON(http.StatusOK, extract)
}

func makeTransaction(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req TransactionBody

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid fields",
		})
		return
	}

	client, ok := database.GetClientInfoCache(id)
	if !ok {
		// user not found
		if !database.DBClient.FindUser(id) {
			ctx.Status(http.StatusNotFound)
			return
		}

		// cache doesnt exists
		transactions, err := database.DBClient.GetAllUserTransactions(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		client = database.CalculateCache(id, transactions)
	}

	// false is invalid transaction because balance is lower than limit
	if !utils.CanMakeTransaction(req.Type, req.Value, client.Balance, client.Limit) {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	transaction, err := database.DBClient.MakeTransaction(client.LastTransactionUUID, id, req.Value, client.Limit, req.Type, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// last saved transaction is not in fact the last transaction, get transactions after last date and calculate cache
	if transaction == nil {
		transactions, err := database.DBClient.GetTransactionsAfterDate(client.UserID, client.LastTransactionDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		client = database.CalculateCache(id, transactions)

		transaction, err = database.DBClient.MakeTransaction(client.LastTransactionUUID, id, req.Value, client.Limit, req.Type, req.Description)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	database.CalculateCache(id, []*database.Transaction{transaction})

	ctx.JSON(http.StatusOK, transaction)
}
