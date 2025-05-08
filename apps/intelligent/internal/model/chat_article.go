package model

import "gorm.io/gorm"

type ChatArticle struct {
	gorm.Model
	ArticleId      int64          `gorm:"column:article_id"`
	Content        string         `gorm:"column:content"`
	ProcessLevel   int64          `gorm:"column:process_level"`
	ProcessAction  ProcessActions `gorm:"column:process_action"`
	ProcessContent string         `gorm:"column:process_content"`
}

func (c *ChatArticle) TableName() string {
	return "chat_article"
}
