package extract

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getExtract(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func CreateRoutes(rt *gin.RouterGroup) {
	rt.GET("/:id/extrato", getExtract)
}
