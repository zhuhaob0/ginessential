package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	mysql := InitDB()
	defer mysql.Close()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	r.POST("api/auth/register", func(c *gin.Context) {
		// 获取参数
		name := c.PostForm("name")
		password := c.PostForm("password")
		telephone := c.PostForm("telephone")

		// 数据验证
		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号必须是11位",
			})
			return
		}

		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能少于6位",
			})
			return
		}

		// 如果名称没有传，给一个10位随机字符串
		if len(name) == 0 {
			name = RandomString(10)
		}

		// 判断手机号是否存在
		if isTelephoneExit(mysql, telephone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户已存在！",
			})
			return
		}
		// 创建用户
		newUser := User{
			Name: name, Telephone: telephone, Password: password,
		}
		mysql.Create(&newUser)

		// 返回结果
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求路径或方法错误",
		})
	})

	panic(r.Run(":8080"))
}

func RandomString(n int) string {
	var letters = []byte("kjagfalihgsikufgvbblacgvbsoljhdjguiwghvbcbmcmlahqwphv")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "127.0.0.1"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := ""
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", username, password, host, port, database, charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
