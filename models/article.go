package models

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title     string         `gorm:"type:varchar(200);not null;index" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Preview   string         `gorm:"type:varchar(500);not null" json:"preview"`
	Likes     int            `gorm:"default:0;index" json:"likes"`
	AuthorID  uint           `gorm:"index" json:"author_id,omitempty"`
	Status    int            `gorm:"default:1;comment:1=published,0=draft" json:"status"`
}

type ArticleQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Title    string `form:"title"`
	Status   int    `form:"status"`
}

func (Article) TableName() string {
	return "articles"
}

func (article *Article) IncrementLikes(db *gorm.DB) error {
	result := db.Model(article).UpdateColumn("likes", gorm.Expr("likes + ?", 1))
	return result.Error
}
