package payload

import (
	"fmt"
)

type Response struct {
	Code      int    `json:"errcode"`
	Msg       string `json:"errmsg,omitempty"`
	Success   bool   `json:"success,omitempty"`
	RequestId string `json:"request_id,omitempty"`
	TraceId   string `json:"requestId,omitempty"` // 新版的api使用TraceId

	// 调用结果
	Result bool `json:"result,omitempty"`
}

// Unmarshalled 统一检查返回异常异常
type Unmarshalled interface {
	CheckError() error
}

func (res *Response) CheckError() (err error) {
	if !res.Ok() {
		err = fmt.Errorf("dingtalk api code:%d, msg: %s", res.Code, res.Msg)
	}
	return err
}

func (res *Response) Ok() bool {
	return res.Code == 0
}
