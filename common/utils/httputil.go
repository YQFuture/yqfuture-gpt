package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPPostAndParseJSON 发起HTTP POST请求，并解析返回的JSON数据
func HTTPPostAndParseJSON(url string, requestData interface{}, responseData interface{}) error {
	// 将请求数据转换为JSON格式
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %v", err)
	}

	// 发起HTTP POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// 解析JSON响应数据
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	return nil
}

// HTTPGetAndParseJSON 发送 HTTP GET 请求并解析返回的 JSON 数据
func HTTPGetAndParseJSON(url string, responseData interface{}) error {
	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET request failed: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP GET request failed with status: %s", resp.Status)
	}

	// 读取并解析 JSON 响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body failed: %v", err)
	}

	err = json.Unmarshal(body, responseData)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	return nil
}
