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

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	client, ok := database.GetClientInfoCache(id)
	if !ok {
		client = getClientInfoAndCache(id)
		if client == nil {
			ctx.Status(http.StatusNotFound)
			return
		}
	} else if !database.DBClient.IsLastTransactionUUID(client.LastTransactionUUID) {
		getClientInfoAndCache(id)
	}

	extract := database.NewExtract(client.Balance, time.Now(), client.Limit, client.LastTransactions)

	ctx.JSON(http.StatusOK, extract)
}

func makeTransaction(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.ParseInt(param, 10, 64)
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
		client = getClientInfoAndCache(id)
		if client == nil {
			ctx.Status(http.StatusNotFound)
			return
		}
	}

	// false is invalid transaction because balance is lower than limit
	if !utils.CanMakeTransaction(req.Type, req.Value, client.Balance, client.Limit) {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	transaction := database.DBClient.MakeTransaction(client.LastTransactionUUID, id, req.Value, req.Type, req.Description)

	// last saved transaction is not in fact the last transaction, get transactions after last date and calculate cache
	if transaction == nil {
		client = getClientInfoAndCache(id)
		if client == nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		transaction = database.DBClient.MakeTransaction(client.LastTransactionUUID, id, req.Value, req.Type, req.Description)
	}

	database.CalculateCache(id, transaction)

	ctx.JSON(http.StatusOK, transaction)
}

func getClientInfoAndCache(id int64) *database.ClientInfo {
	client := database.DBClient.GetClientInfo(id)
	if client == nil {
		return nil
	}

	database.SetClientInfoCache(client)
	return client
}
