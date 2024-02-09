package middlewares

import "github.com/gin-gonic/gin"

func Validate(c *gin.Context) {
	c.Next()
}