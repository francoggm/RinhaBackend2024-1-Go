package client

import (
	"crebito/database"
	"crebito/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	Value       int64  `json:"valor" binding:"required"`
	Type        string `json:"tipo" binding:"required"`
	Description string `json:"descricao" binding:"required"`
}

func CreateRoutes(c *gin.RouterGroup) {
	c.GET("/:id/extrato", getExtract)
	c.POST("/:id/transacoes", makeTransaction)
}

func getExtract(ctx *gin.Context) {
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

	var req Transaction

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid fields",
		})
		return
	}

	client, ok := database.GetClientCache(id)
	if !ok {
		// cache doesnt exists, try to get user from db
		transactions, err := database.DBClient.GetAllUserTransactions(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// user not found
		if len(transactions) < 1 {
			ctx.Status(http.StatusNotFound)
			return
		}
		client = database.CalculateCache(id, transactions)
	}

	// invalid transaction because balance is lower than limit
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

	// last saved transaction is not in fact the last transaction, get transactions after last saved uuid and calculate cache
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
