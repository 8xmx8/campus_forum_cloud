package config

import "campus_forum_cloud/common"

type Config struct {
	MysqlDSN     string                 `json:",optional"`
	CommentChain string                 `json:",optional"`
	ArticleChain string                 `json:",optional"`
	Chat         *ChatConfig            `json:",optional"`
	Redis        *common.RedisClientOpt `json:",optional"`
}

type ChatConfig struct {
	Model   string `json:",optional"`
	BaseUrl string `json:",optional"`
}
