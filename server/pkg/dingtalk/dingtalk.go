package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"godocms/pkg/dingtalk/payload"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DingTalk struct {
	// 企业内部应用对应:AgentId，如果是应套件:SuiteId
	//id int

	// 企业内部应用对应:AppKey，套件对应:SuiteKey
	key string

	// 企业内部对应:AppSecret，套件对应:SuiteSecret
	secret string

	// isv 钉钉开放平台会向应用的回调URL推送的suite_ticket（约5个小时推送一次）
	ticket string

	// 授权企业的id
	corpid string

	// 在开发者后台的基本信息 > 开发信息（旧版）页面获取微应用管理后台SSOSecret
	ssoSecret string

	client *http.Client
	cache  string
}

type OptionFunc func(*DingTalk)

func WithTicket(ticket string) OptionFunc {
	return func(dt *DingTalk) {
		dt.ticket = ticket
	}
}

func WithCorpid(corpid string) OptionFunc {
	return func(dt *DingTalk) {
		dt.corpid = corpid
	}
}

func WithSSOSecret(secret string) OptionFunc {
	return func(dt *DingTalk) {
		dt.ssoSecret = secret
	}
}

func (ding *DingTalk) validate() error {
	if ding.key == "" {
		return errors.New("key is empty")
	}
	if ding.secret == "" {
		return errors.New("secret is empty")
	}
	return nil
}

func NewClient(key, secret string, opts ...OptionFunc) (ding *DingTalk, err error) {
	ding = &DingTalk{key: key, secret: secret}

	for _, opt := range opts {
		opt(ding)
	}

	if err = ding.validate(); err != nil {
		return nil, err
	}

	ding.client = &http.Client{Timeout: 10 * time.Second}

	return
}

func isNewApi(path string) bool {
	return strings.HasPrefix(path, "/v1.0/")
}

func (ding *DingTalk) Request(method, path string, query url.Values, body interface{}, data payload.Unmarshalled) (err error) {
	if body != nil {
		if err := ding.validate(); err != nil {
			return err
		}
	}

	if query == nil {
		query = url.Values{}
	}

	if query.Get("access_token") == "" && path != GetTokenKey && path != CorpAccessToken &&
		path != SuiteAccessToken && path != GetAuthInfo && path != GetAgentKey &&
		path != ActivateSuiteKey && path != GetSSOTokenKey && path != GetUnactiveCorpKey &&
		path != ReauthCorpKey && path != GetCorpPermanentCodeKey && path != GetUserAccessToken {

		var token string
		var err error
		if token, err = ding.GetAccessToken(); err != nil {
			return err
		}
		// set token
		query.Set("access_token", token)
	}

	return ding.httpRequest(method, path, query, body, data)
}

func (ding *DingTalk) httpRequest(method, path string, query url.Values, body interface{}, response payload.Unmarshalled) (err error) {
	var uri *url.URL
	var token string

	newApi := isNewApi(path)
	if newApi {
		token = query.Get("access_token")
		query.Del("access_token")
		uri, _ = url.Parse(NewApi + path)
	} else {
		uri, _ = url.Parse(Api + path)
	}
	uri.RawQuery = query.Encode()

	// 处理请求体
	var reqBody io.Reader
	var contentType string
	if body != nil {
		switch b := body.(type) {
		case string:
			reqBody = strings.NewReader(b)
			contentType = "text/plain"
		case url.Values:
			reqBody = strings.NewReader(b.Encode())
			contentType = "application/x-www-form-urlencoded"
		default:
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			reqBody = bytes.NewReader(jsonBody)
			contentType = "application/json; charset=utf-8"
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, uri.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	if newApi {
		req.Header.Set("x-acs-dingtalk-access-token", token)
	}

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 处理响应
	switch resp.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		return response.CheckError()
	// case http.StatusBadRequest:
	// case http.StatusNotFound:
	// case http.StatusInternalServerError:
	default:
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		return fmt.Errorf("dingtalk server error, unexpected status code: %v, res: %v", resp.StatusCode, result)
	}
}
