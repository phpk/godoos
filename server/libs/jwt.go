package libs

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const secretKey = "godocms2!f" // 使用环境变量或配置文件管理密钥
const expireJwtTimeHours = 8760
const effectTime = time.Duration(expireJwtTimeHours) * time.Hour

// UserClaims 用于存储用户信息，以便生成令牌。
type UserClaims struct {
	ID int64 `json:"id"`
}

// TokenOutTime 返回令牌的过期时间。
func TokenOutTime(claims *UserClaims) int64 {
	return time.Now().Add(effectTime).Unix()
}

// GenerateToken 生成新的JWT令牌。
func GenerateToken(claims *UserClaims) (string, error) {
	claimsStandard := jwt.StandardClaims{
		ExpiresAt: TokenOutTime(claims),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  claims.ID,
		"exp": claimsStandard.ExpiresAt,
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("签名令牌失败: %w", err)
	}

	base58CheckToken, err := EncodeToken(signedToken)
	if err != nil {
		return "", fmt.Errorf("编码令牌失败: %w", err)
	}
	return base58CheckToken, nil
}

// ParseToken 解析和验证JWT令牌。
func ParseToken(base58CheckTokenString string) (*UserClaims, error) {
	token, err := DecodeToken(base58CheckTokenString)
	if err != nil {
		return nil, fmt.Errorf("解码令牌失败: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("无效的令牌")
	}
	id, ok := claims["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("无效的令牌ID")
	}
	return &UserClaims{ID: int64(id)}, nil
}
func Ok(ctx *gin.Context, msg string, data interface{}) {
	token := ctx.Request.Header.Get("Authorization")
	var err error
	if token != "" {
		token, err = Refresh(token)
		if err != nil {
			token = ""
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": 1,
		"message": msg,
		"data":    data,
		"token":   token,
		"time":    time.Now().Unix(),
	})
}

// Refresh 刷新JWT令牌的过期时间。
func Refresh(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("解析令牌失败: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("无效的令牌")
	}
	id, ok := claims["id"].(float64)
	if !ok {
		return "", fmt.Errorf("无效的令牌ID")
	}
	userClaims := &UserClaims{ID: int64(id)}
	claims["exp"] = time.Now().Add(effectTime).Unix()
	token.Claims = claims
	return GenerateToken(userClaims)
}

// EncodeToken 压缩并编码令牌。
func EncodeToken(token string) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(token)); err != nil {
		return "", fmt.Errorf("GZIP写入错误: %w", err)
	}
	if err := gz.Close(); err != nil {
		return "", fmt.Errorf("GZIP关闭错误: %w", err)
	}

	base64UrlSafeToken := base64.URLEncoding.EncodeToString(b.Bytes())
	return base58.CheckEncode([]byte(base64UrlSafeToken), 0x00), nil
}

// DecodeToken 解码并解压令牌。
func DecodeToken(base58CheckTokenString string) (string, error) {
	decodedData, _, err := base58.CheckDecode(base58CheckTokenString)
	if err != nil {
		return "", fmt.Errorf("Base58解码错误: %w", err)
	}

	base64UrlSafeData, err := base64.URLEncoding.DecodeString(string(decodedData))
	if err != nil {
		return "", fmt.Errorf("Base64解码错误: %w", err)
	}

	b := bytes.NewReader(base64UrlSafeData)
	gz, err := gzip.NewReader(b)
	if err != nil {
		return "", fmt.Errorf("GZIP读取器错误: %w", err)
	}
	defer gz.Close()

	var b2 bytes.Buffer
	if _, err := b2.ReadFrom(gz); err != nil {
		return "", fmt.Errorf("GZIP读取错误: %w", err)
	}

	return b2.String(), nil
}
