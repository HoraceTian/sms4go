package sms4go

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// HttpClient 是一个包装 http.Client 的结构体
type HttpClient struct {
	client *http.Client
}

// 单例实例
var (
	instance *HttpClient
	once     sync.Once
)

// NewHttpClient 创建一个新的或返回现有的 HttpClient 实例
func NewHttpClient(client *http.Client) *HttpClient {
	once.Do(func() {
		instance = &HttpClient{
			client: client,
		}
	})
	return instance
}

func GetHttpClient() *HttpClient {
	return instance
}

// PostJson 发送POST JSON请求
func (hc *HttpClient) PostJson(urlStr string, headers map[string]string, payload *TencentPayload) (map[string]interface{}, error) {
	// 1. 处理参数
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("[sms4go] |- failed to marshal request payload")
	}

	// 2. 创建请求
	req, err := http.NewRequest("POST", "https://"+urlStr, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// 3. 处理请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 4. 发起访问
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)

	// 5. 读取结果
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 6. 处理响应
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	// 7. 补偿
	if response, ok := result["Response"].(map[string]interface{}); ok {
		if errMap, ok := response["Error"].(map[string]interface{}); ok {
			if errMsg, ok := errMap["Message"].(string); ok {
				return nil, errors.New(errMsg)
			}
		}
	}

	return result, nil
}

// PostForm 发送POST表单请求
func (hc *HttpClient) PostForm(urlStr string, headers map[string]string, body map[string]string) (map[string]interface{}, error) {
	// 1. 处理参数
	formData := url.Values{}
	for key, value := range body {
		formData.Set(key, value)
	}

	// 2. 创建请求
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	// 3. 创建请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 4. 发起访问
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)

	// 5. 读取结果
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 6. 处理响应
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// PostUrl 发送带URL参数的POST请求
func (hc *HttpClient) PostUrl(baseUrl string, headers map[string]string, params map[string]string) (map[string]interface{}, error) {

	// 1. 处理参数
	urlWithParams := baseUrl + "?" + buildQuery(params)
	req, err := http.NewRequest("POST", urlWithParams, nil)
	if err != nil {
		return nil, err
	}

	// 2. 处理请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 3. 发起访问
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)

	// 4. 读取结果
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 5. 处理响应
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// buildQuery 构建URL查询字符串
func buildQuery(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	return values.Encode()
}

// SafeSleep 线程睡眠
func SafeSleep(retryInterval int) {
	time.Sleep(time.Duration(retryInterval) * time.Second)
}

// closeBody 安全关闭响应体
func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
