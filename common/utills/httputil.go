package utills

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
