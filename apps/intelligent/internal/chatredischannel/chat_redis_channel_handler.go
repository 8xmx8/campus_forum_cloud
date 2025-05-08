package chatredischannel

import (
	"campus_forum_cloud/apps/intelligent/internal/chatredischannel/logic"
	"campus_forum_cloud/apps/intelligent/internal/svc"
	"context"
	"github.com/redis/go-redis/v9"
)

type ChatRedisChannelHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatRedisChannelHandler(ctx context.Context, svcCtx *svc.ServiceContext) *ChatRedisChannelHandler {
	return &ChatRedisChannelHandler{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (c *ChatRedisChannelHandler) Hander(msg *redis.Message) {
	switch msg.Channel {
	case c.svcCtx.Config.CommentChain:
		commentLogic := logic.NewChatCommentLogic(c.ctx, c.svcCtx)
		commentLogic.ChatComment(msg)
	case c.svcCtx.Config.ArticleChain:
		articleLogic := logic.NewChatContentLogic(c.ctx, c.svcCtx)
		articleLogic.ChatContent(msg)
	}
}
