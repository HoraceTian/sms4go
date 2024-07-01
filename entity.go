package sms4go

// TencentPayload 腾讯Payload
type TencentPayload struct {
	PhoneNumberSet   []string `json:"PhoneNumberSet"`
	SignName         string   `json:"SignName"`
	SmsSdkAppId      string   `json:"SmsSdkAppId"`
	TemplateId       string   `json:"TemplateId"`
	TemplateParamSet []string `json:"TemplateParamSet"`
}
