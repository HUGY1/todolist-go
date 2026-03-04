package models

import "time"

type CreateUser struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RemoveUser struct {
	UserID string `json:"userId" binding:"required"`
}

type UpdateUser struct {
	UserID   uint    `json:"userId" binding:"required"`
	UserName *string `json:"userName"`
	Mobile   *string `json:"mobile" `
	Gender   int8    `json:"gender" `
	Avatar   *string `json:"avatar"`
}

type User struct {
	UserID    uint      `json:"userId" gorm:"primaryKey;autoIncrement;comment:用户ID"`
	UserName  string    `json:"userName" gorm:"type:varchar(50);not null;comment:用户名"`
	Mobile    string    `json:"mobile" gorm:"type:varchar(11);uniqueIndex;comment:手机号"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null;comment:密码" `
	Gender    int8      `json:"gender" gorm:"type:tinyint;default:0;comment:性别 0-未知 1-男 2-女"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255);comment:头像"`
	Status    int8      `json:"status" gorm:"type:tinyint;default:0;comment:状态"`
	CreatedAt time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt time.Time `json:"deletedAt" gorm:"comment:删除时间"`
}

func (User) TableName() string {
	return "users"
}

// 性别常量
const (
	GenderUnknown = 0 // 未知
	GenderMale    = 1 // 男
	GenderFemale  = 2 // 女
)
