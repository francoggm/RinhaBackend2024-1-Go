package server

import (
	extract "crebito/controllers/exctract"
	"crebito/controllers/transaction"
	"crebito/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func configureRoutes(r *gin.Engine) {
	clients := r.Group("/clientes")
	clients.Use(middlewares.Validate)
	{
		transaction.CreateRoutes(clients)
		extract.CreateRoutes(clients)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, "Invalid route")
	})
}

func Run(mode string, port string) error {
	gin.SetMode(mode)
	r := gin.Default()

	configureRoutes(r)

	return r.Run(":" + port)
}
