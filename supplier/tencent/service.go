package tencent

import (
	"github.com/panjf2000/ants/v2"
	"log"
	"net/http"
	"sms4go"
	"sms4go/infra"
	"strings"
	"time"
)

type Blender struct {
	parent     sms4go.BaseBlender
	httpClient *infra.HttpClient
}

func (b *Blender) GetConfigId() string {
	return b.parent.ConfigId
}

func (b *Blender) GetSupplier() string {
	return sms4go.Tencent
}

// Helper function to process parameters
func processParameters(message string) map[string]string {
	split := strings.Split(message, sms4go.ParamSeparate)
	paramMap := make(map[string]string, len(split))
	for _, param := range split {
		paramMap[param] = param
	}
	return paramMap
}

// Helper function to create SmsResponse
func (b *Blender) createSmsResponse(phones []string, params []string, templateId string) *sms4go.SmsResponse {
	timestamp := time.Now().Unix()
	cfg := b.parent.Config.(*Config)
	signature, err := generateSignature(cfg, templateId, params, phones, timestamp)
	if err != nil {
		log.Fatalf("[sms4go] |- generate tencent signature failed: %v", err)
		return nil
	}

	headersMap := generateHeadersMap(signature, cfg.action, cfg.version, cfg.territory, cfg.requestUrl, timestamp)
	requestBody := generateRequestBody(phones, cfg.SDKAppId, cfg.Signature, templateId, params)

	resp, err := b.httpClient.PostJson(cfg.requestUrl, headersMap, requestBody)
	if err != nil {
		return sms4go.NewResp(cfg.GetConfigId(), false, err)
	}

	result := b.getResponse(resp)
	if result.Success {
		return sms4go.NewResp(cfg.GetConfigId(), true, result)
	}
	return sms4go.NewFailResp(cfg.GetConfigId(), err)
}

// SendMessage sends a message with parameters
func (b *Blender) SendMessage(phone, message string) *sms4go.SmsResponse {
	return b.SendMessageWithParams(phone, processParameters(message))
}

// SendMessageWithParams sends a message with given parameters
func (b *Blender) SendMessageWithParams(phone string, params map[string]string) *sms4go.SmsResponse {
	templateId := b.parent.Config.(*Config).TemplateId
	return b.SendMessageWithParamsAndTemplate(phone, templateId, params)
}

// SendMessageWithParamsAndTemplate sends a message with parameters and template ID
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
	return b.createSmsResponse(phones, paramList, templateId)
}

// Async message sending helper function
func (b *Blender) sendMessageAsync(sendFunc func() *sms4go.SmsResponse, callback sms4go.Callback) *sms4go.SmsResponse {
	if callback == nil {
		if err := b.parent.RoutinePool.Submit(func() { sendFunc() }); err != nil {
			log.Fatalf("[sms4go] |- submit task failed: %v", err)
		}
		return sms4go.NewResp(b.parent.ConfigId, true, nil)
	}

	responseChan := make(chan *sms4go.SmsResponse, 1)
	if err := b.parent.RoutinePool.Submit(func() {
		response := sendFunc()
		responseChan <- response
		close(responseChan)
	}); err != nil {
		log.Fatalf("[sms4go] |- submit task failed: %v", err)
	}

	go func() {
		response := <-responseChan
		callback(response)
	}()

	return sms4go.NewResp(b.parent.ConfigId, true, nil)
}

// SendMessageAsync Async message sending methods
func (b *Blender) SendMessageAsync(phone, message string, callback sms4go.Callback) *sms4go.SmsResponse {
	return b.sendMessageAsync(func() *sms4go.SmsResponse {
		return b.SendMessage(phone, message)
	}, callback)
}

func (b *Blender) SendMessageWithParamsAsync(phone string, params map[string]string, callback sms4go.Callback) *sms4go.SmsResponse {
	return b.sendMessageAsync(func() *sms4go.SmsResponse {
		return b.SendMessageWithParams(phone, params)
	}, callback)
}

func (b *Blender) SendMessageWithParamsAndTemplateAsync(phone, templateId string, params map[string]string, callback sms4go.Callback) *sms4go.SmsResponse {
	return b.sendMessageAsync(func() *sms4go.SmsResponse {
		return b.SendMessageWithParamsAndTemplate(phone, templateId, params)
	}, callback)
}

// DelayMessage Delay message sending methods
func (b *Blender) DelayMessage(phone, message string, delay time.Duration) {
	b.parent.DelayQueue.AddTask(func() {
		b.SendMessage(phone, message)
	}, delay)
}

func (b *Blender) DelayMessageWithParams(phone string, params map[string]string, delay time.Duration) {
	b.parent.DelayQueue.AddTask(func() {
		b.SendMessageWithParams(phone, params)
	}, delay)
}

func (b *Blender) DelayMessageWithParamsAndTemplate(phone, templateId string, params map[string]string, delay time.Duration) {
	b.parent.DelayQueue.AddTask(func() {
		b.SendMessageWithParamsAndTemplate(phone, templateId, params)
	}, delay)
}

// MassTexting methods
func (b *Blender) MassTexting(phones []string, message string) *sms4go.SmsResponse {
	return b.MassTextingWithParams(phones, processParameters(message))
}

func (b *Blender) MassTextingWithParams(phones []string, params map[string]string) *sms4go.SmsResponse {
	templateId := b.parent.Config.(*Config).TemplateId
	return b.MassTextingWithParamsAndTemplate(phones, templateId, params)
}

func (b *Blender) MassTextingWithParamsAndTemplate(phones []string, templateId string, params map[string]string) *sms4go.SmsResponse {
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
	for index, phone := range phones {
		phones[index] = infra.AddPrefixIfNot(phone, "+86")
	}
	return b.createSmsResponse(phones, paramList, templateId)
}

func (b *Blender) SetHttpClient(client *http.Client) {
	b.httpClient = infra.NewHttpClient(client)
}

func (b *Blender) SetRoutinePool(pool *ants.Pool) {
	b.parent.RoutinePool = pool
}

func (b *Blender) SetDelayQueue(queue *sms4go.DelayQueue) {
	b.parent.DelayQueue = queue
}

func (b *Blender) getResponse(resMap map[string]interface{}) *sms4go.SmsResponse {
	var smsResponse = &sms4go.SmsResponse{}

	response, ok := resMap["Response"].(map[string]interface{})
	if !ok {
		log.Fatalf("[sms4go] |- Response field not found or not a valid object")
		return nil
	}

	errorStr, _ := response["Error"].(string)
	smsResponse.Success = strings.TrimSpace(errorStr) == ""

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

	smsResponse.Data = resMap
	smsResponse.ConfigId = b.parent.ConfigId
	return smsResponse
}
