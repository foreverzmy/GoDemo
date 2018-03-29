package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

// Request 请求结构体
type Request struct {
	URL     string
	Headers map[string]string
	proxy   string
}

// GetHTML 获取 html 页面
func (r Request) GetHTML() []byte {
	// proxy, err := url.Parse(r.proxy)
	// Panic(err)

	client := &http.Client{
		// 设置代理
		// Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
		// 超时时间
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", r.URL, nil)
	Panic(err)

	// 请求里添加 header
	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	Panic(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	Panic(err)

	return body
}
