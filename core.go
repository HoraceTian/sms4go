package sms4go

import (
	"github.com/panjf2000/ants/v2"
	"log"
	"net/http"
	"time"
)

// ISmsBlender 是短信 Blender 的接口，定义了短信服务的基本操作
type ISmsBlender interface {
	// SetRoutinePool 设置协程池
	SetRoutinePool(pool *ants.Pool)

	// SetHttpClient 设置基础 Http 客户端
	SetHttpClient(client *http.Client)

	// SetDelayQueue 设置延时队列
	SetDelayQueue(queue *DelayQueue)

	// GetConfigId 获取配置 ID
	GetConfigId() string

	// GetSupplier 获取供应商信息
	GetSupplier() string

	// SendMessage 发送短信并返回响应
	SendMessage(phone, message string) *SmsResponse

	// SendMessageWithParams 发送短信并使用给定参数返回响应
	SendMessageWithParams(phone string, params map[string]string) *SmsResponse

	// SendMessageWithParamsAndTemplate 发送短信并使用给定参数与指定模版返回响应
	SendMessageWithParamsAndTemplate(phone, templateId string, params map[string]string) *SmsResponse

	// SendMessageAsync 异步发送短信并返回响应
	SendMessageAsync(phone, message string, callback Callback) *SmsResponse

	// SendMessageWithParamsAsync 异步发送短信并使用给定参数返回响应
	SendMessageWithParamsAsync(phone string, params map[string]string, callback Callback) *SmsResponse

	// SendMessageWithParamsAndTemplateAsync 异步发送具有模板 ID 和参数的短信并返回响应
	SendMessageWithParamsAndTemplateAsync(phone, templateId string, params map[string]string, callback Callback) *SmsResponse

	// DelayMessage 延时消息
	DelayMessage(phone, message string, delay time.Duration)

	// DelayMessageWithParams 延迟消息带参数
	DelayMessageWithParams(phone string, params map[string]string, delay time.Duration)

	// DelayMessageWithParamsAndTemplate 延迟消息带参数模版
	DelayMessageWithParamsAndTemplate(phone, templateId string, params map[string]string, delay time.Duration)

	// MassTexting 群发消息
	MassTexting(phones []string, message string) *SmsResponse

	// MassTextingWithParams 群发消息带参数
	MassTextingWithParams(phones []string, params map[string]string) *SmsResponse

	// MassTextingWithParamsAndTemplate 群发消息带参数模版
	MassTextingWithParamsAndTemplate(phones []string, templateId string, params map[string]string) *SmsResponse
}

// BaseBlender 共有抽象
type BaseBlender struct {
	ConfigId    string
	Config      SupplierConfig
	RoutinePool *ants.Pool
	DelayQueue  *DelayQueue
}

func (b *BaseBlender) sendMessageAsync(sendFunc func() *SmsResponse, callback Callback) *SmsResponse {
	// 如果没有提供回调函数，直接提交任务到协程池执行
	if callback == nil {
		if err := b.RoutinePool.Submit(func() { sendFunc() }); err != nil {
			log.Fatalf("[sms4go] |- submit task failed: %v", err)
		}
		return NewResp(b.ConfigId, true, nil)
	}

	// 创建一个 channel 来传递 SmsResponse
	responseChan := make(chan *SmsResponse, 1)
	if err := b.RoutinePool.Submit(func() {
		response := sendFunc()
		responseChan <- response
		close(responseChan)
	}); err != nil {
		log.Fatalf("[sms4go] |- submit task failed: %v", err)
	}

	// 启动一个 goroutine 来处理回调
	go func() {
		response := <-responseChan
		callback(response)
	}()

	return NewResp(b.ConfigId, true, nil)
}
