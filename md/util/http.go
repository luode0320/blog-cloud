package util

import (
	"bytes"
	"md/middleware"
	"net/http"
)

func RefreshDir() {
	// 定义请求体，这里假设POST请求不需要携带额外的数据
	requestBody := []byte{}

	// 创建HTTP客户端
	client := &http.Client{}

	// 构建请求
	req, err := http.NewRequest("POST", "http://0.0.0.0:4000/refresh-dir", bytes.NewBuffer(requestBody))
	if err != nil {
		middleware.Log.Warnf("创建请求失败: %s", err)
	}
	// 设置请求头，模拟浏览器请求（可选，根据需求调整）
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		middleware.Log.Warnf("发送请求失败: %s", err)
	}
	if resp == nil {
		return
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
}
