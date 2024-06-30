package tencent

import (
	"fmt"
	"sms4go"
	"sms4go/infra"
	"strings"
)

type Blender struct {
	parent sms4go.BaseBlender
}

func (b Blender) GetConfigId() string {
	return b.parent.ConfigId
}

func (b Blender) GetSupplier() string {
	return infra.Tencent
}

func (b Blender) SendMessage(phone, message string) *sms4go.SmsResponse {
	// 1. 处理参数
	split := strings.Split(message, infra.ParamSeparate)
	paramMap := make(map[string]string, len(split))
	for _, param := range split {
		paramMap[param] = param
	}

	// 2. 发送
	return b.SendMessageWithParams(phone, paramMap)
}

func (b Blender) SendMessageWithParams(phone string, params map[string]string) *sms4go.SmsResponse {
	// 1. 提取模版号
	templateId := b.parent.Config.(*Config).TemplateId

	// 2. 发送
	return b.SendMessageWithParamsAndTemplate(phone, templateId, params)
}

func (b Blender) SendMessageAsync(phone, message string, callback sms4go.Callback) *sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b Blender) SendMessageWithParamsAsync(phone string, params map[string]string, callback sms4go.Callback) *sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b Blender) SendMessageWithTemplateAsync(phone, templateId string, params map[string]string, callback sms4go.Callback) *sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b Blender) SendMessageWithParamsAndTemplate(phone, templateId string, params map[string]string) *sms4go.SmsResponse {
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

func (b Blender) getSmsResponse(phones []string, params []string, templateId string) *sms4go.SmsResponse {
	fmt.Println("调用 GetSmsResponse 了，", phones, params, templateId)
	return &sms4go.SmsResponse{}
}
