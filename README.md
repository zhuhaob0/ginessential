## 基本功能
请求注册用户，密码通过**bcrypt**加密，存储用户信息至数据库，**jwt**生成并发放token，登录时**bcript**验证密码正确性，携带token请求个人信息。

文章分类，创建、删除、修改、查看文章类别标签操作。

文章管理，发布、删除、修改、查看文章操作

注：开发时postman测试请求参数保存在**ginEssential.postman_collection.json**文件中

## 系统软件环境
```text
系统：ubuntu-22.04.2
数据库：mysql-8.0.32
GO版本：1.18.3
```

注：项目使用第三方包版本如**go.mod**所示

## 代码目录
```text
.
├── common
├── config
├── controller
├── dto
├── ginessential
├── ginEssential.postman_collection.json
├── go.mod
├── go.sum
├── main.go
├── middleware
├── model
├── README.md
├── repository
├── response
├── routes.go
├── util
└── vo
```
**common**：封装对外公共调用的方法

**config**：配置文件

**controller**：控制器，处理不同的路由请求

**dto**：全称 Data Transfer Object，做表示层，展示给用户的数据对象

**middleware**：中间件

**model**：数据库中表的结构体模型

**repository**：将操作数据库表的函数统一封装在一个结构体里，便于controller层调用

**response**：给前端返回的统一格式

**util**：工具包

**vo**：全称 View Object，绑定前端传来的数据

## 启动命令

`go build && ./ginessential`
