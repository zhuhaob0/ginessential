package controller

import (
	"net/http"
	"oceanlearn/ginessential/common"
	"oceanlearn/ginessential/model"
	"oceanlearn/ginessential/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Register(c *gin.Context) {
	mysql := common.GetDB()
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
		name = util.RandomString(10)
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
	newUser := model.User{
		Name: name, Telephone: telephone, Password: password,
	}
	mysql.Create(&newUser)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"msg": "注册成功",
	})
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
