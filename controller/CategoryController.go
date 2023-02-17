package controller

import (
	"oceanlearn/ginessential/model"
	"oceanlearn/ginessential/repository"
	"oceanlearn/ginessential/response"
	"oceanlearn/ginessential/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	repository := repository.NewCategoryRepository()
	repository.DB.AutoMigrate(model.Category{})

	return CategoryController{Repository: repository}
}

// 添加分类
func (c CategoryController) Create(ctx *gin.Context) {
	// 绑定body的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误："+err.Error())
		return
	}
	category, err := c.Repository.Create(requestCategory.Name)
	if err != nil {
		panic(err)
	}
	response.Success(ctx, gin.H{"category": category}, "创建成功")
}

// 修改分类
func (c CategoryController) Update(ctx *gin.Context) {
	// 绑定body中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	// 获取path中的参数
	categoryID, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		response.Fail(ctx, nil, "请求的id不是数字")
		return
	}
	updateCategory, err := c.Repository.SelectById(categoryID)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}
	// 更新分类
	category, err := c.Repository.Update(*updateCategory, requestCategory.Name)
	if err != nil {
		panic(err)
	}
	response.Success(ctx, gin.H{"category": category}, "修改成功")
}

// 删除分类
func (c CategoryController) Delete(ctx *gin.Context) {
	// 获取path中的参数
	categoryID, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		response.Fail(ctx, nil, "请求的id不是数字")
		return
	}
	if err := c.Repository.DeleteById(categoryID); err != nil {
		response.Fail(ctx, nil, "删除失败，请重试")
		return
	}
	response.Success(ctx, nil, "删除成功")
}

// 查询分类
func (c CategoryController) Show(ctx *gin.Context) {
	// 获取path中的参数
	categoryID, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		response.Fail(ctx, nil, "请求的id不是数字")
		return
	}
	category, err := c.Repository.SelectById(categoryID)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}
	response.Success(ctx, gin.H{"category": category}, "分类存在")
}
