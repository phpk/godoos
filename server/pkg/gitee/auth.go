package gitee

import (
	"context"
	"encoding/json"
	"errors"
	"godocms/common"
	"godocms/pkg/gitee/payload"
	"log/slog"

	"golang.org/x/oauth2"
)

// OAuthConfig 配置 OAuth2 客户端
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Endpoint     oauth2.Endpoint
	Scopes       []string
}

// GiteeOAuth 结构体用于处理 Gitee OAuth2 相关的操作
type GiteeOAuth struct {
	config *oauth2.Config
}

// NewGithubOAuth 创建一个新的 GithubOAuth 实例
func NewGithubOAuth(oauthConfig OAuthConfig) *GiteeOAuth {
	if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" || oauthConfig.RedirectURL == "" {
		oauthConfig = newDefaultOAuthConfig()
	}

	return &GiteeOAuth{
		config: &oauth2.Config{
			ClientID:     oauthConfig.ClientID,
			ClientSecret: oauthConfig.ClientSecret,
			Endpoint:     oauthConfig.Endpoint,
			RedirectURL:  oauthConfig.RedirectURL,
			Scopes:       oauthConfig.Scopes,
		},
	}
}

func newDefaultOAuthConfig() OAuthConfig {
	return OAuthConfig{
		ClientID:     common.LoginConf.Gitee.ClientID,
		ClientSecret: common.LoginConf.Gitee.ClientSecret,
		RedirectURL:  common.LoginConf.Gitee.RedirectURL,
		Scopes:       common.LoginConf.Gitee.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://gitee.com/oauth/authorize",
			TokenURL: "https://gitee.com/oauth/token",
		},
	}
}

// GetAccessToken 使用 code 获取 access_token
func (g *GiteeOAuth) GetAccessToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		slog.Error("获取access_token失败", "err", err)
		return nil, err
	}
	return token, nil
}

// GetUserInfo 使用 access_token 获取用户信息
func (g *GiteeOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (payload.User, error) {
	var giteeUser payload.User
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://gitee.com/api/v5/user")
	if err != nil {
		slog.Error("获取用户信息失败", "err", err)
		return giteeUser, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&giteeUser); err != nil {
		slog.Error("解析用户信息失败", "err", err)
		return giteeUser, err
	}

	if giteeUser.ID == 0 {
		return giteeUser, errors.New("github user id is empty")
	}

	return giteeUser, nil
}

func (g *GiteeOAuth) GetGiteeOAuthAuthUrl(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
