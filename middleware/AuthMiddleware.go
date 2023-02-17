package middleware

import (
	"net/http"
	"oceanlearn/ginessential/common"
	"oceanlearn/ginessential/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取authorization header
		tokenString := c.GetHeader("Authorization")
		// validate token format
		//oauth2.0规定的,Authorization的字符串开头必须要有"Bearer "
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token为空或不合法"})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token可能过期或其他错误"})
			c.Abort()
			return
		}
		// 验证通过后获取claim中的userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 用户
		if userId == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未查询到用户"})
			c.Abort()
			return
		}
		// 用户存在，将user的信息写入上下文
		c.Set("user", user)
		c.Next()
	}
}
