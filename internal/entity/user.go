package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Email     string    `json:"email" gorm:"size:100;uniqueIndex"`
	Password  string    `json:"-" gorm:"size:100;not null"`
	Nickname  string    `json:"nickname" gorm:"size:50"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Status    int       `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}