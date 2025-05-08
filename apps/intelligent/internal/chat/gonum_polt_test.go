package chat

import (
	"campus_forum_cloud/apps/intelligent/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommentMessage(t *testing.T) {
	client := NewOllamaClient("http://localhost:11434", "llama3.2:1b")
	message, err := client.CommentMessage("你写的太垃圾了")
	assert.NoError(t, err)
	t.Log(message.Messages.Content)
	char := utils.GetFirstChar(message.Messages.Content)
	t.Log(char)

}
