package logic

import (
	"campus_forum_cloud/apps/intelligent/internal/chat"
	"campus_forum_cloud/apps/intelligent/internal/model"
	"campus_forum_cloud/apps/intelligent/internal/svc"
	"campus_forum_cloud/apps/intelligent/internal/utils"
	"campus_forum_cloud/common"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatContentLogic {
	return &ChatContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ChatContentLogic) ChatContent(msg *redis.Message) {
	contentIdStr := msg.Payload
	contentId := cast.ToInt64(contentIdStr)
	if contentId == 0 {
		return
	}
	logx.Infof("内容id为 :%d\n", contentId)
	contentData, err := l.svcCtx.DAO.GetContentByID(contentId)
	if err != nil {
		logx.Errorf("获取内容失败 :%v\n", err)
		return
	}
	if contentData.Content == "" {
		logx.Errorf("内容内容为空 :%v\n", err)
		return
	}
	ok := common.Validate(l.svcCtx.SnFilter, contentData.Content)
	if ok {
		err = l.svcCtx.DAO.UpdateContentNoOllama(contentId)
		if err != nil {
			logx.Errorf("更新帖子失败: %v\n", err)
		}
		return
	}
	client := chat.NewOllamaClient(l.svcCtx.Config.Chat.BaseUrl, l.svcCtx.Config.Chat.Model)
	chatResponse, err := client.CommentMessage(contentData.Content)
	if err != nil {
		logx.Errorf("大模型审核帖子失败 :%v\n", err)
		return
	}
	level := utils.GetFirstChar(chatResponse.Messages.Content)
	err = l.svcCtx.DAO.UpdateArticleContent(&model.ChatArticle{
		ArticleId:      contentId,
		Content:        contentData.Content,
		ProcessLevel:   level,
		ProcessAction:  model.GetProcessAction(level),
		ProcessContent: chatResponse.Messages.Content,
	})
	if err != nil {
		logx.Errorf("更新帖子失败 :%v\n", err)
		return
	}
}
