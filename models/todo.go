// Package models 定义应用的数据模型
package models

import (
	"time"
)

// Todo 待办事项模型
// 对应数据库中的 list 表（按需求使用 'list' 作为表名）
// GORM 会自动处理表的创建和字段映射
type Todo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`      // 主键，自增
	Value       string    `gorm:"type:varchar(500);not null" json:"value"` // 待办事项内容
	IsCompleted bool      `gorm:"default:false" json:"isCompleted"`        // 是否完成，默认 false
	CreatedAt   time.Time `json:"createdAt"`                               // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`
	FileId      *uint     `json:"fileId" gorm:"default:null"`              // 关联文件，可为空
	File        *File     `json:"file,omitempty" gorm:"foreignKey:FileId"` // 更新时间
}

// TableName 指定表名为 list
// 实现 Tabler 接口，使 GORM 使用 'list' 作为表名而非默认的 todos
func (Todo) TableName() string {
	return "list"
}

// AddTodoRequest 添加待办事项的请求体结构
// 用于 POST /api/add-todo 接口的 JSON 解析
type AddTodoRequest struct {
	Value       string `json:"value" binding:"required"` // 待办内容，必填
	IsCompleted *bool  `json:"isCompleted"`              // 是否完成，可选，默认 false
}

type GetTodoRequest struct {
	StartIndex uint    `json:"startIndex" form:"startIndex" binding:"omitempty,min=0"`     // 开始索引，可选
	PageSize   uint    `json:"pageSize" form:"pageSize" binding:"omitempty,min=1,max=100"` // 每页大小，可选
	SearchKey  *string `json:"searchKey" form:"searchKey"`                                 // 搜索关键字
}

type ListDataWithPagination[T any] struct {
	Items      []T  `json:"items"`
	TotalItems int  `json:"totalItems"`
	PageSize   uint `json:"pageSize"`
}
