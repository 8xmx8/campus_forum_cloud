package model

import (
	"gorm.io/gorm"
	"time"
)

const CommentContentIllegal = "该评论已违规"

// CampusComment 表示数据库中的 campus_comment 表。
type CampusComment struct {
	// CommentID 是评论的主键。
	CommentID int64 `gorm:"primaryKey;column:comment_id;type:bigint;not null" json:"comment_id"`

	// ParentID 是父评论的ID。默认为0。
	ParentID int64 `gorm:"column:parent_id;type:bigint;default:0" json:"parent_id"`

	// OneLevelID 是该评论所属的一级评论的ID。默认为-1。
	OneLevelID int64 `gorm:"column:one_level_id;type:bigint;default:-1" json:"one_level_id"`

	// UserID 是发表评论的用户的ID。
	UserID int64 `gorm:"column:user_id;type:bigint" json:"user_id"`

	// ToUserID 是评论指向的用户的ID。默认为-1。
	ToUserID int64 `gorm:"column:to_user_id;type:bigint;default:-1" json:"to_user_id"`

	// ContentID 是评论关联的内容的ID。
	ContentID int64 `gorm:"column:content_id;type:bigint" json:"content_id"`

	// CoContent 是评论的内容，最大长度为200个字符。
	CoContent string `gorm:"column:co_content;type:varchar(200);charset:utf8mb4;collate:utf8mb4_general_ci" json:"co_content"`

	// IP 是用户发表评论时的IP地址。
	IP string `gorm:"column:ip;type:varchar(64);charset:utf8mb4;collate:utf8mb4_general_ci" json:"ip"`

	// Address 是用户发表评论时的地址。
	Address string `gorm:"column:address;type:varchar(100);charset:utf8mb4;collate:utf8mb4_general_ci" json:"address"`

	// CreateTime 是评论创建的时间戳。默认为当前时间。
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP" json:"create_time"`

	// CreateUser 是创建评论的用户的ID。
	CreateUser int64 `gorm:"column:create_user;type:bigint" json:"create_user"`

	// UpdateTime 是评论最后更新的时间戳。更新时自动设置为当前时间。
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"update_time"`

	// UpdateUser 是最后更新评论的用户的ID。
	UpdateUser int64 `gorm:"column:update_user;type:bigint" json:"update_user"`
}

// TableName 设置 CampusComment 模型的表名。
func (c *CampusComment) TableName() string {
	return "campus_comment"
}

// 根据评论ID获取评论
func (d *DAO) GetCommentByID(commentID int64) (*CampusComment, error) {
	var comment CampusComment
	err := d.DB.Where("comment_id = ?", commentID).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// 修改评论内容,添加大模型审核
func (d *DAO) UpdateCommentContent(chatComment *ChatComment) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&ChatComment{}).Create(&chatComment).Error
		if err != nil {
			return err
		}
		if chatComment.ProcessLevel > 1 {
			err = tx.Model(&CampusComment{}).
				Where("comment_id = ?", chatComment.CommentId).
				Update("co_content", CommentContentIllegal).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func (d *DAO) UpdateCommentNoOllama(id int64) error {
	return d.DB.Model(&CampusComment{}).
		Where("comment_id = ?", id).
		Update("co_content", CommentContentIllegal).Error
}
