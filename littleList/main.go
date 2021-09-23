package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

// to do model

type ToDoItem struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

var db *gorm.DB

func initDb() (err error) {
	// 连接数据库
	dsn := "root:a662449aa@(xlsf.xyz:3306)/test?charset=UTF8&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}

func main() {
	//连接数据库
	err := initDb()
	if err != nil {
		return
	}

	// 创建表
	db.AutoMigrate(&ToDoItem{})

	// 定义默认路由
	engine := gin.Default()

	// 加载资源
	engine.LoadHTMLGlob("./templates/*")

	// 加载静态资源
	engine.Static("/static", "static")

	// 路由函数，get请求index
	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// 待办事项
	v1Group := engine.Group("/v1")
	{
		// 添加
		// 前端页面填写待办事项，点击提交，请求方式为post
		v1Group.POST("/todo", func(context *gin.Context) {
			// 从请求中把数据拿出来
			var item ToDoItem
			err := context.ShouldBind(&item)
			if err != nil {
				fmt.Println("read error...")
				return
			}
			//存入数据库,返回响应
			if err = db.Create(&item).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"code": 2000,
					"msg":  "success",
					"data": item,
				})
			}
		})

		// 查看所有代办事项

		v1Group.GET("/todo", func(context *gin.Context) {
			// 查询数据表里的所有数据
			var toDoItem []ToDoItem
			if err = db.Find(&toDoItem).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, toDoItem)
			}
		})

		// 查看某个待办事项

		v1Group.GET("/todo/:id", func(context *gin.Context) {

		})


		// 修改

		v1Group.PUT("/todo/:id", func(context *gin.Context) {
			id, _ := context.Params.Get("id")
			var toDoItem ToDoItem
			if err = db.Where("id = ?", id).First(&toDoItem).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{"error": err})
				return
			}
			context.BindJSON(&toDoItem) // 绑定属性
			if err = db.Save(&toDoItem).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{"error": err})
			} else {
				context.JSON(http.StatusOK, toDoItem)
			}
		})

		// 删除

		v1Group.DELETE("/todo/:id", func(context *gin.Context) {
			id, ok := context.Params.Get("id")
			if !ok {
				context.JSON(http.StatusOK, gin.H{"error": "无效id"})
				return
			}
			if err = db.Where("id = ?", id).Delete(&ToDoItem{}).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{"error": err})
			} else {
				context.JSON(http.StatusOK, gin.H{"status": "deleted"})
			}
		})

	}

	//启动服务器
	err = engine.Run(":1234")
	if err != nil {
		fmt.Println("start service failed...", err)
	}
}
