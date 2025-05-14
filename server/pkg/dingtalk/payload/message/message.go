package message

type message struct {
	// 消息类型
	MsgType string `json:"msgtype" validate:"required,oneof=text image voice file link oa markdown action_card feedCard"`
}

// Message 消息结构
type Message interface {
	// MessageType 消息类型
	MessageType() string

	String() string
}

// Response 发送消息返回
type Response struct {
	Code      int    `json:"errcode"`          // code
	Msg       string `json:"errmsg,omitempty"` // msg
	Success   bool   `json:"success,omitempty"`
	RequestId string `json:"request_id,omitempty"`

	MessageId string `json:"messageId"`
}
