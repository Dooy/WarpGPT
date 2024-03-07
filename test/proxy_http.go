package test

import (
	"WarpGPT/pkg/logger"
	"fmt"
	"io"
	"net/http"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("有请求过来了！ ", r.URL.String())
	//fmt.Println("有请求过来了！!")
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 从原始请求中复制HTTP头到新请求
	req.Header = r.Header
	// 创建一个HTTP客户端
	client := &http.Client{}
	// 发送请求到目标服务器
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	// 复制目标服务器响应到原始请求中
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func StartHttpProxy(server string) {
	http.HandleFunc("/", ProxyHandler)
	fmt.Println("代理服务器已经运行,监听端口:", server)
	err := http.ListenAndServe(server, nil) //":8080"
	if err != nil {
		fmt.Println("代理服务器启动失败!")
	}
}
