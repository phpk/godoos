package alibaba

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"godocms/common"
	"io"
	"sort"

	"golang.org/x/exp/maps"

	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	httpMethod   string
	canonicalUri string
	host         string
	xAcsAction   string
	xAcsVersion  string
	headers      map[string]string
	body         []byte
	queryParam   map[string]interface{}
}

func NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion string) *Request {
	req := &Request{
		httpMethod:   httpMethod,
		canonicalUri: canonicalUri,
		host:         host,
		xAcsAction:   xAcsAction,
		xAcsVersion:  xAcsVersion,
		headers:      make(map[string]string),
		queryParam:   make(map[string]interface{}),
	}
	req.headers["host"] = host
	req.headers["x-acs-action"] = xAcsAction
	req.headers["x-acs-version"] = xAcsVersion
	req.headers["x-acs-date"] = time.Now().UTC().Format(time.RFC3339)
	req.headers["x-acs-signature-nonce"] = uuid.New().String()
	return req
}

// os.Getenv()表示从环境变量中获取AccessKey ID和AccessKey Secret。
var (
	AccessKeyId     = common.LoginConf.Phone.AliyunSms.AccessKeyId
	AccessKeySecret = common.LoginConf.Phone.AliyunSms.AccessKeySecret
	ALGORITHM       = "ACS3-HMAC-SHA256"
)

// 签名示例，您需要根据实际情况替换main方法中的示例参数。
// ROA接口和RPC接口只有canonicalUri取值逻辑是完全不同，其余内容都是相似的。
// 通过API元数据获取请求方法（methods）、请求参数名称（name）、请求参数类型（type）、请求参数位置（in），并将参数封装到SignatureRequest中。
// 1. 请求参数在元数据中显示"in":"query"，通过queryParam传参。
// 2. 请求参数在元数据中显示"in": "body"，通过body传参。
// 3. 请求参数在元数据中显示"in": "formData"，通过body传参。
func signExample() {
	// RPC接口请求示例一：请求参数"in":"query"
	httpMethod := "POST"                   // 请求方式，大部分RPC接口同时支持POST和GET，此处以POST为例
	canonicalUri := "/"                    // RPC接口无资源路径，故使用正斜杠（/）作为CanonicalURI
	host := "ecs.cn-hangzhou.aliyuncs.com" // 云产品服务接入点
	xAcsAction := "SendSms"                // API名称
	xAcsVersion := "2014-05-26"            // API版本号
	req := NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion)
	// DescribeInstanceStatus请求参数如下：
	// RegionId在元数据中显示的类型是String，"in":"query"，必填
	req.queryParam["RegionId"] = "cn-hangzhou"
	// InstanceId的在元数据中显示的类型是array，"in":"query"，非必填
	instanceIds := []interface{}{"i-bp10igfmnyttXXXXXXXX", "i-bp1incuofvzxXXXXXXXX", "i-bp1incuofvzxXXXXXXXX"}
	req.queryParam["InstanceId"] = instanceIds

	// // RPC接口请求示例二：请求参数"in":"body"
	// httpMethod := "POST"
	// canonicalUri := "/"
	// host := "ocr-api.cn-hangzhou.aliyuncs.com"
	// xAcsAction := "RecognizeGeneral"
	// xAcsVersion := "2021-07-07"
	// req := NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion)
	// // 读取文件内容
	// filePath := "D:\\test.png"
	// bytes, err := os.ReadFile(filePath)
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }
	// req.body = bytes
	// req.headers["content-type"] = "application/octet-stream"

	// // RPC接口请求示例三：请求参数"in": "formData"
	// httpMethod := "POST"
	// canonicalUri := "/"
	// host := "mt.aliyuncs.com"
	// xAcsAction := "TranslateGeneral"
	// xAcsVersion := "2018-10-12"
	// req := NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion)
	// // TranslateGeneral请求参数如下：
	// // Context在元数据中显示的类型是String，"in":"query"，非必填
	// req.queryParam["Context"] = "早上"
	// // FormatType、SourceLanguage、TargetLanguage等参数，在元数据中显示"in":"formData"
	// body := make(map[string]interface{})
	// body["FormatType"] = "text"
	// body["SourceLanguage"] = "zh"
	// body["TargetLanguage"] = "en"
	// body["SourceText"] = "你好"
	// body["Scene"] = "general"
	// str := formDataToString(body)
	// req.body = []byte(*str)
	// req.headers["content-type"] = "application/x-www-form-urlencoded"

	// // ROA接口POST请求
	// httpMethod := "POST"
	// canonicalUri := "/clusters"
	// host := "cs.cn-beijing.aliyuncs.com"
	// xAcsAction := "CreateCluster"
	// xAcsVersion := "2015-12-15"
	// req := NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion)
	// // 封装请求参数，请求参数在元数据中显示"in": "body"，表示参数放在body中
	// body := make(map[string]interface{})
	// body["name"] = "testDemo"
	// body["region_id"] = "cn-beijing"
	// body["cluster_type"] = "ExternalKubernetes"
	// body["vpcid"] = "vpc-2zeou1uod4ylaXXXXXXXX"
	// body["container_cidr"] = "10.0.0.0/8"
	// body["service_cidr"] = "172.16.1.0/20"
	// body["security_group_id"] = "sg-2ze1a0rlgeo7XXXXXXXX"
	// vswitch_ids := []interface{}{"vsw-2zei30dhfldu8XXXXXXXX"}
	// body["vswitch_ids"] = vswitch_ids
	// jsonBytes, err := json.Marshal(body)
	// if err != nil {
	// 	fmt.Println("Error marshaling to JSON:", err)
	// 	return
	// }
	// req.body = []byte(jsonBytes)
	// req.headers["content-type"] = "application/json; charset=utf-8"

	// // ROA接口GET请求
	// httpMethod := "GET"
	// // canonicalUri如果存在path参数，需要对path参数encode，percentCode({path参数})
	// canonicalUri := "/clusters/" + percentCode("c558c166928f9446dae400d106e124f66") + "/resources"
	// host := "cs.cn-beijing.aliyuncs.com"
	// xAcsAction := "DescribeClusterResources"
	// xAcsVersion := "2015-12-15"
	// req := NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion)
	// req.queryParam["with_addon_resources"] = "true"

	// // ROA接口DELETE请求
	// httpMethod := "DELETE"
	// // canonicalUri如果存在path参数，需要对path参数encode，percentCode({path参数})
	// canonicalUri := "/clusters/" + percentCode("c558c166928f9446dae400d106e124f66")
	// host := "cs.cn-beijing.aliyuncs.com"
	// xAcsAction := "DeleteCluster"
	// xAcsVersion := "2015-12-15"
	// req := NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion)

	// 签名过程
	getAuthorization(req)
	// 调用API
	_, error := callAPI(req)
	if error != nil {
		println(error.Error())
	}
}

func callAPI(req *Request) ([]byte, error) {
	urlStr := "https://" + req.host + req.canonicalUri
	q := url.Values{}
	keys := maps.Keys(req.queryParam)
	sort.Strings(keys)
	for _, k := range keys {
		v := req.queryParam[k]
		q.Set(k, fmt.Sprintf("%v", v))
	}
	urlStr += "?" + q.Encode()
	fmt.Println(urlStr)

	httpReq, err := http.NewRequest(req.httpMethod, urlStr, strings.NewReader(string(req.body)))
	if err != nil {
		return nil, err
	}

	for key, value := range req.headers {
		httpReq.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var respBuffer bytes.Buffer
	_, err = io.Copy(&respBuffer, resp.Body)
	if err != nil {
		return nil, err
	}
	respBytes := respBuffer.Bytes()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d, %s", resp.StatusCode, respBytes)
	}

	// fmt.Println(string(respBytes))
	return respBytes, nil
}

func getAuthorization(req *Request) {
	// 处理queryParam中参数值为List、Map类型的参数，将参数平铺
	newQueryParams := make(map[string]interface{})
	processObject(newQueryParams, "", req.queryParam)
	req.queryParam = newQueryParams
	// 步骤 1：拼接规范请求串
	canonicalQueryString := ""
	keys := maps.Keys(req.queryParam)
	sort.Strings(keys)
	for _, k := range keys {
		v := req.queryParam[k]
		canonicalQueryString += percentCode(url.QueryEscape(k)) + "=" + percentCode(url.QueryEscape(fmt.Sprintf("%v", v))) + "&"
	}
	canonicalQueryString = strings.TrimSuffix(canonicalQueryString, "&")
	// fmt.Printf("canonicalQueryString========>%s\n", canonicalQueryString)

	var bodyContent []byte
	if req.body == nil {
		bodyContent = []byte("")
	} else {
		bodyContent = req.body
	}
	hashedRequestPayload := sha256Hex(bodyContent)
	req.headers["x-acs-content-sha256"] = hashedRequestPayload

	canonicalHeaders := ""
	signedHeaders := ""
	HeadersKeys := maps.Keys(req.headers)
	sort.Strings(HeadersKeys)
	for _, k := range HeadersKeys {
		lowerKey := strings.ToLower(k)
		if lowerKey == "host" || strings.HasPrefix(lowerKey, "x-acs-") || lowerKey == "content-type" {
			canonicalHeaders += lowerKey + ":" + req.headers[k] + "\n"
			signedHeaders += lowerKey + ";"
		}
	}
	signedHeaders = strings.TrimSuffix(signedHeaders, ";")

	canonicalRequest := req.httpMethod + "\n" + req.canonicalUri + "\n" + canonicalQueryString + "\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + hashedRequestPayload
	// fmt.Printf("canonicalRequest========>\n%s\n", canonicalRequest)

	// 步骤 2：拼接待签名字符串
	hashedCanonicalRequest := sha256Hex([]byte(canonicalRequest))
	stringToSign := ALGORITHM + "\n" + hashedCanonicalRequest
	// fmt.Printf("stringToSign========>\n%s\n", stringToSign)

	// 步骤 3：计算签名
	byteData, err := hmac256([]byte(AccessKeySecret), stringToSign)
	if err != nil {
		fmt.Println(err)
	}
	signature := strings.ToLower(hex.EncodeToString(byteData))

	// 步骤 4：拼接Authorization
	authorization := ALGORITHM + " Credential=" + AccessKeyId + ",SignedHeaders=" + signedHeaders + ",Signature=" + signature
	// fmt.Printf("authorization========>%s\n", authorization)
	req.headers["Authorization"] = authorization
}

func hmac256(key []byte, toSignString string) ([]byte, error) {
	// 实例化HMAC-SHA256哈希
	h := hmac.New(sha256.New, key)
	// 写入待签名的字符串
	_, err := h.Write([]byte(toSignString))
	if err != nil {
		return nil, err
	}
	// 计算签名并返回
	return h.Sum(nil), nil
}

func sha256Hex(byteArray []byte) string {
	// 实例化SHA-256哈希函数
	hash := sha256.New()
	// 将字符串写入哈希函数
	_, _ = hash.Write(byteArray)
	// 计算SHA-256哈希值并转换为小写的十六进制字符串
	hexString := hex.EncodeToString(hash.Sum(nil))

	return hexString
}

func percentCode(str string) string {
	// 替换特定的编码字符
	str = strings.ReplaceAll(str, "+", "%20")
	str = strings.ReplaceAll(str, "*", "%2A")
	str = strings.ReplaceAll(str, "%7E", "~")
	return str
}

func formDataToString(formData map[string]interface{}) *string {
	tmp := make(map[string]interface{})
	processObject(tmp, "", formData)
	res := ""
	urlEncoder := url.Values{}
	for key, value := range tmp {
		v := fmt.Sprintf("%v", value)
		urlEncoder.Add(key, v)
	}
	res = urlEncoder.Encode()
	return &res
}

// processObject 递归处理对象，将复杂对象（如Map和List）展开为平面的键值对
func processObject(mapResult map[string]interface{}, key string, value interface{}) {
	if value == nil {
		return
	}

	switch v := value.(type) {
	case []interface{}:
		for i, item := range v {
			processObject(mapResult, fmt.Sprintf("%s.%d", key, i+1), item)
		}
	case map[string]interface{}:
		for subKey, subValue := range v {
			processObject(mapResult, fmt.Sprintf("%s.%s", key, subKey), subValue)
		}
	default:
		if strings.HasPrefix(key, ".") {
			key = key[1:]
		}
		if b, ok := v.([]byte); ok {
			mapResult[key] = string(b)
		} else {
			mapResult[key] = fmt.Sprintf("%v", v)
		}
	}
}
