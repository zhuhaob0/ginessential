package controller

import (
	"oceanlearn/ginessential/common"
	"oceanlearn/ginessential/model"
	"oceanlearn/ginessential/response"
	"oceanlearn/ginessential/vo"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}

// 创建文章
func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	//绑定body参数，数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误："+err.Error())
		return
	}
	// 获取登录用户 user
	user, _ := ctx.Get("user")
	// 创建post
	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}
	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
	}
	response.Success(ctx, nil, "创建成功")
}

// 更新文章
func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	//绑定body参数，数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误："+err.Error())
		return
	}
	// 获取path中的 id
	postId := ctx.Params.ByName("id")
	var post model.Post
	if p.DB.Preload("Category").Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	// 获取登录用户 user，判断是否为文章作者
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "您不是文章的作者，请勿非法操作")
		return
	}
	// 更新文章
	if err := p.DB.Preload("Category").Model(&post).Update(requestPost).Error; err != nil {
		response.Fail(ctx, nil, "更新失败："+err.Error())
		return
	}
	response.Success(ctx, gin.H{"post": post}, "更新成功")
}

// 查看文章
func (p PostController) Show(ctx *gin.Context) {
	// 获取path中的 id
	postId := ctx.Params.ByName("id")
	var post model.Post
	if p.DB.Preload("Category").Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	response.Success(ctx, gin.H{"post": post}, "查询成功")
}

// 删除文章
func (p PostController) Delete(ctx *gin.Context) {
	// 获取path中的 id
	postId := ctx.Params.ByName("id")
	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	// 获取登录用户 user，判断是否为文章作者
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "您不是文章的作者，请勿非法操作")
		return
	}
	if err := p.DB.Delete(&post).Error; err != nil {
		response.Fail(ctx, nil, "删除失败："+err.Error())
		return
	}
	response.Success(ctx, gin.H{"post": post}, "删除成功")
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	var pageNum, pageSize int
	var err error
	if pageNum, err = strconv.Atoi(ctx.DefaultQuery("pageNum", "1")); err != nil {
		response.Fail(ctx, nil, "pageNum参数必须是数字")
		return
	}
	if pageSize, err = strconv.Atoi(ctx.DefaultQuery("pageSize", "20")); err != nil {
		response.Fail(ctx, nil, "pageSize参数必须是数字")
		return
	}
	// 分页
	var posts []model.Post
	p.DB.Preload("Category").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// 前端渲染分页需要知道的总数
	var total int
	p.DB.Model(model.Post{}).Count(&total)

	response.Success(ctx, gin.H{"data": posts, "total": total}, "分页成功")
}
