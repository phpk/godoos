package dingtalk

import (
	"godocms/pkg/dingtalk/payload"
	"net/http"
)

// SendCorpConvMessage 发送工作通知
func (ding *DingTalk) SendCorpConvMessage(req *payload.CorpConvMessage) (rsp payload.CorpConvMessageResponse,
	err error) {
	return rsp, ding.Request(http.MethodPost, SendCorpConversationMessageKey, nil, req, &rsp)
}
