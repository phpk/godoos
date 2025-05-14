package message

import (
	"encoding/json"
)

// 钉钉消息结构体
type text struct {
	Content string `json:"content" validate:"required"`
}

// 文本消息
type textMessage struct {
	message
	text `json:"text" validate:"required"`
}

func (t *textMessage) String() string {
	str, _ := json.Marshal(t)
	return string(str)
}

func (t *textMessage) MessageType() string {
	return "text"
}

// NewTextMessage 文本对象
func NewTextMessage(context string) *textMessage {
	msg := &textMessage{}
	msg.MsgType = msg.MessageType()
	msg.text = text{Content: context}
	return msg
}
