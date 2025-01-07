package search

import (
	"net"
	"net/http"
	"time"
)

var tr *http.Transport

func init() {

	tr = &http.Transport{
		// 设置代理
		MaxIdleConns: 100,
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(time.Second * 6)) //设置发送接受数据超时
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
