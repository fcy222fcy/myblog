package entity

import (
	"time"
)

// BaseEntity 基础实体
type BaseEntity struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
