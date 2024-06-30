package tencent

import (
	"fmt"
	"log"
	"net/http"
	"sms4go"
	"sms4go/infra"
	"strings"
	"time"
)

var retry = 0

type Blender struct {
	parent     sms4go.BaseBlender
	httpClient *infra.HttpClient
}

func (b *Blender) GetConfigId() string {
	return b.parent.ConfigId
}

func (b *Blender) GetSupplier() string {
	return infra.Tencent
}

func (b *Blender) SendMessage(phone, message string) *sms4go.SmsResponse {
	// 1. 处理参数
	split := strings.Split(message, infra.ParamSeparate)
	paramMap := make(map[string]string, len(split))
	for _, param := range split {
		paramMap[param] = param
	}

	// 2. 发送
	return b.SendMessageWithParams(phone, paramMap)
}

func (b *Blender) SendMessageWithParams(phone string, params map[string]string) *sms4go.SmsResponse {
	// 1. 提取模版号
	templateId := b.parent.Config.(*Config).TemplateId

	// 2. 发送
	return b.SendMessageWithParamsAndTemplate(phone, templateId, params)
}

func (b *Blender) SendMessageAsync(phone, message string, callback sms4go.Callback) *sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b *Blender) SendMessageWithParamsAsync(phone string, params map[string]string, callback sms4go.Callback) *sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b *Blender) SendMessageWithTemplateAsync(phone, templateId string, params map[string]string, callback sms4go.Callback) *sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b *Blender) SendMessageWithParamsAndTemplate(phone, templateId string, params map[string]string) *sms4go.SmsResponse {
	if templateId == "" {
		panic("[sms4go] |- the templateId is empty")
	}
	if params == nil {
		params = make(map[string]string)
	}
	paramList := make([]string, 0, len(params))
	for _, value := range params {
		paramList = append(paramList, value)
	}
	phones := []string{infra.AddPrefixIfNot(phone, "+86")}
	return b.getSmsResponse(phones, paramList, templateId)
}

func (b *Blender) getSmsResponse(phones []string, params []string, templateId string) *sms4go.SmsResponse {
	// 1. 获取时间
	timestamp := time.Now().Unix()

	// 2. 生成签名
	cfg := b.parent.Config.(*Config)
	signature, err := generateSignature(cfg, templateId, params, phones, timestamp)
	fmt.Println(signature)
	if err != nil {
		log.Fatalf("[sms4go] |- generate tencent signature failed: %v", err)
		return nil
	}

	// 3. 处理请求头与请求体
	headersMap := generateHeadersMap(signature, cfg.action, cfg.version, cfg.territory, cfg.requestUrl, timestamp)
	requestBody := generateRequestBody(phones, cfg.SDKAppId, cfg.Signature, templateId, params)

	// 4. 发送请求
	resp, err := b.httpClient.PostJson(cfg.requestUrl, headersMap, requestBody)
	if err != nil {
		return sms4go.NewResp(cfg.GetConfigId(), false, err)
	}

	// 5. 处理响应
	result := b.getResponse(resp)
	if result.Success {
		return sms4go.NewResp(cfg.GetConfigId(), true, result)
	}
	return sms4go.NewFailResp(cfg.GetConfigId(), result)
}

func (b *Blender) getResponse(resMap map[string]interface{}) *sms4go.SmsResponse {
	var smsResponse = &sms4go.SmsResponse{} // 初始化 SmsResponse 指针

	// 获取 Response 字段
	response, ok := resMap["Response"].(map[string]interface{})
	if !ok {
		log.Fatalf("[sms4go] |-  Response field not found or not a valid object")
		return nil
	}

	// 根据 Error 判断是否配置错误
	errorStr, _ := response["Error"].(string)
	smsResponse.Success = strings.TrimSpace(errorStr) == ""

	// 根据 SendStatusSet 判断是否不为Ok
	sendStatusSet, ok := response["SendStatusSet"].([]interface{})
	if ok {
		success := true
		for _, obj := range sendStatusSet {
			if jsonObject, ok := obj.(map[string]interface{}); ok {
				code, _ := jsonObject["Code"].(string)
				if code != "Ok" {
					success = false
					break
				}
			}
		}
		smsResponse.Success = success
	}

	// 设置 SmsResponse 的其他字段
	smsResponse.Data = resMap
	smsResponse.ConfigId = b.parent.ConfigId

	// 返回 SmsResponse 指针
	return smsResponse
}

func (b *Blender) SetHttpClient(client *http.Client) {
	b.httpClient = infra.NewHttpClient(client)
}
