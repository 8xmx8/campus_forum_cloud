package model

import (
	"gorm.io/gorm"
	"time"
)

type CampusContent struct {
	ContentID   int64     `json:"contentId" gorm:"column:content_id;primaryKey"`
	UserID      int64     `json:"userId" gorm:"column:user_id"`
	CategoryID  int64     `json:"categoryId" gorm:"column:category_id"`
	Content     string    `json:"content" gorm:"column:content;size:700"`
	Status      int8      `json:"status" gorm:"column:status"` // 状态：0审核,1正常,2下架,3拒绝（审核不通过）
	Type        int8      `json:"type" gorm:"column:type"`     // 类型：0文字,1图片,2视频
	FileCount   int       `json:"fileCount" gorm:"column:file_count"`
	LoveCount   int       `json:"loveCount" gorm:"column:love_count;default:0"`
	IsAnonymous int8      `json:"isAnonymous" gorm:"column:is_anonymous;default:0"`
	Remark      string    `json:"remark" gorm:"column:remark;size:500"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	CreateUser  *int64    `json:"createUser" gorm:"column:create_user"`
	UpdateTime  time.Time `json:"updateTime" gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdateUser  *int64    `json:"updateUser" gorm:"column:update_user"`
}

func (c *CampusContent) TableName() string {
	return "campus_content"
}

func (d *DAO) GetContentByID(contentID int64) (*CampusContent, error) {
	var content CampusContent
	err := d.DB.Where("content_id = ?", contentID).First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (d *DAO) UpdateArticleContent(chatContent *ChatArticle) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&ChatArticle{}).Create(chatContent).Error
		if err != nil {
			return err
		}
		_UpdateMap := map[string]interface{}{
			"update_time": time.Now(),
			"update_user": 1,
			"status":      1,
		}
		_SqlBase := d.DB.Model(&CampusContent{}).
			Where("content_id = ?", chatContent.ArticleId)
		if chatContent.ProcessLevel > 1 {
			_UpdateMap["status"] = 3
			return _SqlBase.Updates(_UpdateMap).Error

		}
		return _SqlBase.Updates(_UpdateMap).Error
	})
}

func (d *DAO) UpdateContentNoOllama(contentID int64) error {
	return d.DB.Model(&CampusContent{}).
		Where("content_id = ?", contentID).
		Updates(map[string]interface{}{
			"update_time": time.Now(),
			"update_user": 1,
			"status":      3,
		}).Error
}
