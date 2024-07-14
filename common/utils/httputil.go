package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
)

// HTTPPostAndParseJSON 发起HTTP POST请求，并解析返回的JSON数据
func HTTPPostAndParseJSON(url string, requestData interface{}, responseData interface{}) error {
	logx.Infof("发送 HTTP POST 请求url: %s, 请求体: %v", url, requestData)
	// 将请求数据转换为JSON格式
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("解析JSON失败, err: %v", err)
	}

	// 发起HTTP POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("发送POST请求失败, err: %v", err)
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
		return fmt.Errorf("读取响应体失败, err: %v", err)
	}

	// 解析JSON响应数据
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return fmt.Errorf("解析JSON失败, err :%v", err)
	}

	logx.Infof("发送 HTTP POST 响应体: %v", responseData)
	return nil
}

// HTTPGetAndParseJSON 发送 HTTP GET 请求并解析返回的 JSON 数据
func HTTPGetAndParseJSON(url string, responseData interface{}) error {
	logx.Infof("发送 HTTP GET 请求url: %s", url)
	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("发送GET请求失败, err: %v", err)
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
		return fmt.Errorf("读取响应体失败, err: %v", err)
	}

	err = json.Unmarshal(body, responseData)
	if err != nil {
		return fmt.Errorf("解析JSON失败, err: %v", err)
	}

	logx.Infof("发送 HTTP GET 响应体: %v", responseData)
	return nil
}
