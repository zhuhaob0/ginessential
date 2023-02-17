package middleware

import (
	"fmt"
	"oceanlearn/ginessential/response"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(c, nil, fmt.Sprint(err))
			}
		}()
		c.Next()
	}
}
