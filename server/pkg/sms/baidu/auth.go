package baidu

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Auth is auth module of baidu api authenticator
type auth struct {
	AccessKey string
	SecretKey string
}

// getCanonicalTime returns RFC3339 style time string
// current param provides convenience for testing
func getCanonicalTime(current time.Time) (isoStr string) {
	utc, _ := time.LoadLocation("UTC")
	isoStr = current.In(utc).Format(time.RFC3339)
	return
}

func hash(data string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func stringNormalize(str string, encodingSlash bool) string {
	var b strings.Builder
	for _, ch := range str {
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' || ch == '~' || ch == '.' {
			b.WriteRune(rune(ch))
		} else if ch == '/' {
			if encodingSlash {
				b.WriteString("%2F")
			} else {
				b.WriteRune(rune(ch))
			}
		} else {
			// query escape turn space to +, but baidu api turn space to %20
			// issue: https://github.com/golang/go/issues/4013
			b.WriteString(strings.ReplaceAll(url.QueryEscape(string(ch)), "+", "%20"))
		}
	}
	return b.String()
}

func uriCanonicalization(path string) string {
	return path
}

func queryStringCanonicalization(query url.Values) string {
	var queryPatternArray []string
	for k, values := range query {
		if strings.ToLower(k) == "authorization" {
			continue
		}
		if len(values) > 0 {
			for _, v := range values {
				queryPatternArray = append(queryPatternArray, fmt.Sprintf("%s=%s", k, stringNormalize(v, true)))
			}
		} else {
			queryPatternArray = append(queryPatternArray, fmt.Sprintf("%s=", k))
		}
	}
	sort.Strings(queryPatternArray)
	return strings.Join(queryPatternArray, "&")
}

func headersCanonicalization(headers http.Header) (string, string) {
	headersToSign := map[string]bool{"host": true, "content-md5": true, "content-length": true, "content-type": true}
	var canonicalHeaders []string
	var signedHeaders []string
	for k, values := range headers {
		k = strings.ToLower(k)
		if headersToSign[k] || strings.HasPrefix(k, "x-bce-") {
			signedHeaders = append(signedHeaders, k)
			for _, v := range values {
				if v != "" {
					canonicalHeaders = append(canonicalHeaders, fmt.Sprintf("%s:%s", stringNormalize(k, true), stringNormalize(v, true)))
				}
			}
		}
	}
	sort.Strings(canonicalHeaders)
	sort.Strings(signedHeaders)
	return strings.Join(canonicalHeaders, "\n"), strings.Join(signedHeaders, ";")
}

var defaultExpirationInSeconds = 1800

// timestamp int64, expirationInSeconds int, headersToSign []string 这个三个不支持了
func (auth auth) generateAuthorization(httpMethod string, path string, headers http.Header, query url.Values, currentTime time.Time) string {
	rawSessionKey := fmt.Sprintf("bce-auth-v1/%s/%s/%d", auth.AccessKey, getCanonicalTime(currentTime), defaultExpirationInSeconds)
	sessionKey := hash(rawSessionKey, auth.SecretKey)
	canonicalURI := uriCanonicalization(path)
	canonicalQueryString := queryStringCanonicalization(query)
	canonicalHeaders, signedHeaders := headersCanonicalization(headers)
	rawSignature := fmt.Sprintf("%s\n%s\n%s\n%s", httpMethod, canonicalURI, canonicalQueryString, canonicalHeaders)
	signature := hash(rawSignature, sessionKey)
	if len(signedHeaders) > 0 {
		return fmt.Sprintf("%s/%s/%s", rawSessionKey, signedHeaders, signature)
	}
	return fmt.Sprintf("%s//%s", rawSessionKey, signature)
}
