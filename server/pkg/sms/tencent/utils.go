package tencent

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type callbackFunc func(error, *http.Response, string)

type option struct {
	Protocol string            `json:"protocol"`
	Host     string            `json:"host"`
	Path     string            `json:"path"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Body     interface{}       `json:"body"`
}

type tel struct {
	NationCode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}

func getRandom() int {
	return int(math.Round(rand.Float64() * 99999))
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}

func calculateSignature(appKey string, random int, time int64, phoneNumbers []string) string {
	var err error
	h := sha256.New()
	if len(phoneNumbers) != 0 {
		var PhoneNumbers string
		for _, v := range phoneNumbers {
			PhoneNumbers = PhoneNumbers + v + ","
		}
		PhoneNumbers = PhoneNumbers[:len(PhoneNumbers)-1]
		_, err = h.Write([]byte("appkey=" + appKey + "&random=" + strconv.Itoa(random) + "&time=" + strconv.FormatInt(time, 10) + "&mobile=" + PhoneNumbers))
		if err != nil {
			return ""
		}
	} else {
		_, err = h.Write([]byte("appkey=" + appKey + "&random=" + strconv.Itoa(random) + "&time=" + strconv.FormatInt(time, 10)))
		if err != nil {
			return ""
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

func calculateAuth(appKey string, random int, time int64, fileSha1Sum string) string {
	var err error
	h := sha256.New()
	_, err = h.Write([]byte("appkey=" + appKey + "&random=" + strconv.Itoa(random) + "&time=" + strconv.FormatInt(time, 10) + "&content-sha1=" + fileSha1Sum))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func sha1sum(buf []byte) string {
	s := sha1.New()
	s.Write([]byte(buf))
	return hex.EncodeToString(s.Sum(nil))
}

func request(options option, callback callbackFunc) error {
	var err error
	var rawBody []byte
	rawBody, err = json.Marshal(options.Body)
	if err != nil {
		return nil
	}
	body := bytes.NewReader(rawBody)
	url := options.Protocol + `://` + options.Host + options.Path
	var req *http.Request
	req, err = http.NewRequest(options.Method, url, body)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	for k, v := range options.Headers {
		req.Header.Add(k, v)
	}
	var resp *http.Response
	c := http.Client{}
	resp, err = c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if callback != nil {
		callback(err, resp, string(data))
	}
	return err
}
