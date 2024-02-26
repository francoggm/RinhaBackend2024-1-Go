package server

import (
	"crebito/controllers/client"
	"crebito/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func configureRoutes(r *gin.Engine) {
	c := r.Group("/clientes")
	c.Use(middlewares.Validate)
	client.CreateRoutes(c)

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, "Invalid route")
	})
}

func Run(mode string, port string) error {
	gin.SetMode(mode)
	r := gin.Default()

	configureRoutes(r)

	return r.Run("0.0.0.0:" + port)
}
