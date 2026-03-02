// Package handlers 处理 HTTP 请求和响应
package handlers

import (
	"log"
	"net/http"
	"strconv"

	"todolist/config"
	"todolist/models"

	"github.com/gin-gonic/gin"
)

// GetTodos 获取所有待办事项
// 接口: GET /api/get-todo
// 从 list 表中查询并返回所有待办事项
func GetTodos(c *gin.Context) {
	log.Println("[GetTodos] 收到获取待办列表请求")

	var req models.GetTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AddTodo] 请求参数无效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效，startIndex，pageSize 为必填项",
			"error":   err.Error(),
		})
		return
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // 限制最大每页数量
	}

	query := config.GetDB().Model(&models.Todo{})
	if req.SearchKey != nil {
		query = query.Where("value LIKE ?", "%"+*req.SearchKey+"%")
	}
	// query = query.Order("created_at DESC").Offset(int(req.StartIndex)).Limit(int(req.PageSize))
	query = query.Order("created_at DESC")

	var todos []models.Todo
	if err := query.Find(&todos).Error; err != nil {
		log.Printf("[GetTodos] 查询失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询待办事项失败",
			"error":   err.Error(),
		})
		return
	}

	// var nextIndex uint

	if len(todos) > int(req.PageSize) {
		todos = todos[req.StartIndex : req.StartIndex+req.PageSize] // 只返回 pageSize 条数据
	}

	// if len(todos) > 0 {
	// 	nextIndex = todos[len(todos)-1].ID
	// }

	// 若没有数据，返回空数组而非 null
	if todos == nil {
		todos = []models.Todo{}
	}

	log.Printf("[GetTodos] 成功返回 %d 条待办事项", len(todos))

	result := models.ListDataWithPagination[models.Todo]{Items: todos, TotalItems: len(todos), PageSize: req.PageSize}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AddTodo 添加新的待办事项
// 接口: POST /api/add-todo
// 请求体: {"value": string, "isCompleted": boolean(可选)}
// 返回新添加的待办事项，包含自动生成的 id
func AddTodo(c *gin.Context) {
	log.Println("[AddTodo] 收到添加待办请求")

	var req models.AddTodoRequest
	// 使用 ShouldBindJSON 解析并验证请求体
	// binding:"required" 确保 value 字段存在且非空
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AddTodo] 请求参数无效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效，value 为必填项",
			"error":   err.Error(),
		})
		return
	}

	// 处理 isCompleted 默认值
	// 若未传入则默认为 false
	isCompleted := false
	if req.IsCompleted != nil {
		isCompleted = *req.IsCompleted
	}

	todo := models.Todo{
		Value:       req.Value,
		IsCompleted: isCompleted,
	}

	// 插入数据库
	if err := config.GetDB().Create(&todo).Error; err != nil {
		log.Printf("[AddTodo] 插入失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "添加待办事项失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[AddTodo] 成功添加待办事项, id=%d", todo.ID)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    todo,
	})
}

// UpdateTodo 更新待办事项的完成状态（取反）
// 接口: POST /api/update-todo/:id
// 路径参数 id 为待办事项的唯一标识
// 将 isCompleted 值取反后返回更新后的对象
func UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("[UpdateTodo] 收到更新请求, id=%s", idStr)

	// 解析 id 为整数
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("[UpdateTodo] id 格式无效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id 必须是有效的数字",
			"error":   err.Error(),
		})
		return
	}

	var todo models.Todo
	// 先根据 id 查询记录是否存在
	if err := config.GetDB().First(&todo, id).Error; err != nil {
		log.Printf("[UpdateTodo] 未找到 id=%d 的待办事项: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "待办事项不存在",
			"error":   err.Error(),
		})
		return
	}

	// 将 isCompleted 取反
	todo.IsCompleted = !todo.IsCompleted

	// 更新数据库
	if err := config.GetDB().Save(&todo).Error; err != nil {
		log.Printf("[UpdateTodo] 更新失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新待办事项失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[UpdateTodo] 成功更新 id=%d, isCompleted=%v", todo.ID, todo.IsCompleted)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    todo,
	})
}

// DeleteTodo 根据 id 删除待办事项
// 接口: POST /api/del-todo/:id
// 路径参数 id 为待办事项的唯一标识
// 返回删除操作的结果状态
func DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("[DeleteTodo] 收到删除请求, id=%s", idStr)

	// 解析 id 为整数
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("[DeleteTodo] id 格式无效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id 必须是有效的数字",
			"error":   err.Error(),
		})
		return
	}

	// 执行删除操作
	// Delete 会返回影响的行数
	result := config.GetDB().Delete(&models.Todo{}, id)
	if result.Error != nil {
		log.Printf("[DeleteTodo] 删除失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除待办事项失败",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[DeleteTodo] 未找到 id=%d 的待办事项", id)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "待办事项不存在",
		})
		return
	}

	log.Printf("[DeleteTodo] 成功删除 id=%d", id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
		"data": gin.H{
			"id":      id,
			"deleted": true,
		},
	})
}
