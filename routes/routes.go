// Package routes 定义 API 路由和中间件
package routes

import (
	"log"
	"todolist/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回 Gin 引擎实例
// 注册所有 API 路由，遵循 RESTful 设计
func SetupRouter() *gin.Engine {
	// 创建默认的 Gin 引擎（包含 Logger 和 Recovery 中间件）
	r := gin.Default()

	// 全局日志中间件：记录每个请求的基本信息
	r.Use(func(c *gin.Context) {
		log.Printf("[%s] %s - %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	})

	// API 路由组：所有接口都以 /api 为前缀
	api := r.Group("/api")
	{
		// GET /api/get-todo - 查询所有待办事项
		api.POST("/get-todo", handlers.GetTodos)

		// POST /api/add-todo - 添加新的待办事项
		api.POST("/add-todo", handlers.AddTodo)

		// POST /api/update-todo/:id - 根据 id 更新待办事项状态（isCompleted 取反）
		api.POST("/update-todo/:id", handlers.UpdateTodo)

		// POST /api/del-todo/:id - 根据 id 删除待办事项
		api.POST("/del-todo/:id", handlers.DeleteTodo)

		// POST /api/upload - 根据 id 删除待办事项
		api.POST("/upload", handlers.Upload)

		// POST /api/create-user - 创建用户
		api.POST("/create-user", handlers.CreateUser)
		// POST /api/update-user - 更新用户
		api.POST("/update-user", handlers.UpdateUser)
	}

	return r
}
