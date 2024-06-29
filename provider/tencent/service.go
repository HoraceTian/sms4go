package tencent

import (
	"fmt"
	"sms4go"
	"sms4go/provider"
)

type Blender struct {
	parent sms4go.BaseBlender
}

func (b Blender) GetConfigId() string {
	return b.parent.ConfigId
}

func (b Blender) GetSupplier() string {
	return provider.Tencent
}

func (b Blender) SendMessage(phone, message string) sms4go.SmsResponse {
	fmt.Println("哈哈哈哈， 我是腾讯")
	return sms4go.SmsResponse{}
}

func (b Blender) SendMessageWithParams(phone string, params map[string]string) sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b Blender) SendMessageWithTemplate(phone, templateId string, params map[string]string) sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b Blender) SendMessageAsync(phone, message string, callback sms4go.Callback) sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b Blender) SendMessageWithParamsAsync(phone string, params map[string]string, callback sms4go.Callback) sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}

func (b Blender) SendMessageWithTemplateAsync(phone, templateId string, params map[string]string, callback sms4go.Callback) sms4go.SmsResponse {
	//TODO implement me
	panic("implement me")
}
