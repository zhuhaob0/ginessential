package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Post struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserId     uint      `json:"user_id" gorm:"not null"`     // 作者id
	CategoryId uint      `json:"category_id" gorm:"not null"` //文章分类id
	Category   *Category
	Title      string `json:"title" gorm:"type:varchar(50); not null"` //标题
	HeadImg    string `json:"head_img" gorm:"type:varchar(255)"`
	Content    string `json:"content" gorm:"type:text; not null"` // 文章图片储存地址
	CreatedAt  Time   `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt  Time   `json:"updated_at" gorm:"type:timestamp"`
}

func (post *Post) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
