package main

import (
	"oceanlearn/ginessential/common"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysql := common.InitDB()
	defer mysql.Close()

	r := gin.Default()
	r = CollectRoute(r)

	panic(r.Run(":8080"))
}
