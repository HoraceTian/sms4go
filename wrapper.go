package sms4go

// NewSuccessResp 成功响应
func NewSuccessResp(configId string, data interface{}) *SmsResponse {
	return &SmsResponse{
		Success:  true,
		ConfigId: configId,
		Data:     data,
	}
}

// NewResp 任意结果
func NewResp(configId string, resultFlag bool, data interface{}) *SmsResponse {
	return &SmsResponse{
		Success:  resultFlag,
		ConfigId: configId,
		Data:     data,
	}
}

func NewFailResp(configId string, err error) *SmsResponse {
	return &SmsResponse{
		Success:  false,
		ConfigId: configId,
		Data:     err,
	}
}
