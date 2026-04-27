package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Email     string         `gorm:"type:varchar(100);index" json:"email,omitempty"`
	Status    int            `gorm:"default:1;comment:1=active,0=inactive" json:"status"`
}

func (User) TableName() string {
	return "users"
}
