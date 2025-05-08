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

type ChatCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatCommentLogic {
	return &ChatCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ChatCommentLogic) ChatComment(msg *redis.Message) {
	commentIdStr := msg.Payload
	commentId := cast.ToInt64(commentIdStr)
	if commentId == 0 {
		return
	}
	logx.Infof("评论id为 :%d\n", commentId)
	commentData, err := l.svcCtx.DAO.GetCommentByID(commentId)
	if err != nil {
		logx.Errorf("获取评论失败 :%v\n", err)
		return
	}
	if commentData.CoContent == "" {
		logx.Errorf("评论内容为空 :%v\n", err)
		return
	}

	ok := common.Validate(l.svcCtx.SnFilter, commentData.CoContent)

	if ok {
		err = l.svcCtx.DAO.UpdateCommentNoOllama(commentId)
		if err != nil {
			logx.Errorf("更新评论失败 :%v\n", err)
		}

		return
	}

	client := chat.NewOllamaClient(l.svcCtx.Config.Chat.BaseUrl, l.svcCtx.Config.Chat.Model)

	chatResponse, err := client.CommentMessage(commentData.CoContent)
	if err != nil {
		logx.Errorf("大模型审核评论失败 :%v\n", err)
		return
	}
	level := utils.GetFirstChar(chatResponse.Messages.Content)
	err = l.svcCtx.DAO.UpdateCommentContent(&model.ChatComment{
		CommentId:      commentId,
		Content:        commentData.CoContent,
		ProcessLevel:   level,
		ProcessAction:  model.GetProcessAction(level),
		ProcessContent: chatResponse.Messages.Content,
	})
	if err != nil {
		logx.Errorf("更新评论失败 :%v\n", err)
		return
	}
}
