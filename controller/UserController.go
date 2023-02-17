package controller

import (
	"log"
	"net/http"
	"oceanlearn/ginessential/common"
	"oceanlearn/ginessential/dto"
	"oceanlearn/ginessential/model"
	"oceanlearn/ginessential/response"
	"oceanlearn/ginessential/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	mysql := common.GetDB()

	// 使用map 获取请求参数
	// var requestMap = make(map[string]string)
	// json.NewDecoder(c.Request.Body).Decode(&requestMap)

	var requestUser = model.User{}
	// 使用结构体 获取参数
	// json.NewDecoder(c.Request.Body).Decode(&requestUser)
	// 或者使用gin的Bind方法获取参数
	c.Bind(&requestUser)

	// 获取参数，注释为在postman的form-data数据请求时使用的获取参数的方法
	name := requestUser.Name           // c.PostForm("name")
	password := requestUser.Password   // c.PostForm("password")
	telephone := requestUser.Telephone // c.PostForm("telephone")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		return
	}

	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 如果名称没有传，给一个10位随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	// 判断手机号是否存在
	if isTelephoneExit(mysql, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}
	// 创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name: name, Telephone: telephone, Password: string(hashedPassword),
	}
	mysql.Create(&newUser)
	// 发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}

	// 返回结果
	response.Success(c, gin.H{"token": token}, "注册成功")
}

func Login(c *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	c.Bind(&requestUser)
	// 获取参数，注释为在postman的form-data数据请求时使用的获取参数的方法
	password := requestUser.Password   // c.PostForm("password")
	telephone := requestUser.Telephone // c.PostForm("telephone")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		return
	}

	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 判断手机号是否存在

	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}

	// 返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")

}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	response.Response(c, http.StatusOK, 200, gin.H{"user": dto.ToUserDto(user.(model.User))}, "用户基础信息")
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
