package github

import (
	"context"
	"encoding/json"
	"errors"
	"godocms/common"
	"godocms/pkg/github/payload"
	"log/slog"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// OAuthConfig 配置 OAuth2 客户端
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// GithubOAuth 结构体用于处理 GitHub OAuth2 相关的操作
type GithubOAuth struct {
	config *oauth2.Config
}

// NewGithubOAuth 创建一个新的 GithubOAuth 实例
func NewGithubOAuth(oauthConfig OAuthConfig) *GithubOAuth {
	if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" || oauthConfig.RedirectURL == "" {
		oauthConfig = newDefaultOAuthConfig()
	}

	return &GithubOAuth{
		config: &oauth2.Config{
			ClientID:     oauthConfig.ClientID,
			ClientSecret: oauthConfig.ClientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  oauthConfig.RedirectURL,
			Scopes:       oauthConfig.Scopes,
		},
	}
}

func newDefaultOAuthConfig() OAuthConfig {
	return OAuthConfig{
		ClientID:     common.LoginConf.Github.ClientID,
		ClientSecret: common.LoginConf.Github.ClientSecret,
		RedirectURL:  common.LoginConf.Github.RedirectURL,
		Scopes:       common.LoginConf.Github.Scopes,
	}
}

// GetAccessToken 使用 code 获取 access_token
func (g *GithubOAuth) GetAccessToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		slog.Error("获取access_token失败", "err", err)
		return nil, err
	}
	return token, nil
}

// GetUserInfo 使用 access_token 获取用户信息
func (g *GithubOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (payload.User, error) {
	var githubUser payload.User
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		slog.Error("获取用户信息失败", "err", err)
		return githubUser, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		slog.Error("解析用户信息失败", "err", err)
		return githubUser, err
	}

	if githubUser.ID == 0 {
		return githubUser, errors.New("github user id is empty")
	}

	return githubUser, nil
}

func (g *GithubOAuth) GetGithubAuthUrl(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
