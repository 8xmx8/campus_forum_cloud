package model

import "gorm.io/gorm"

type ChatComment struct {
	gorm.Model
	CommentId      int64          `gorm:"column:comment_id"`
	Content        string         `gorm:"column:content"`
	ProcessLevel   int64          `gorm:"column:process_level"`
	ProcessAction  ProcessActions `gorm:"column:process_action"`
	ProcessContent string         `gorm:"column:process_content"`
}

func (c *ChatComment) TableName() string {
	return "chat_comment"
}
