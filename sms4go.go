package sms4go

// NewSmsClient 初始化 Sms 客户端
func NewSmsClient(options ...Option) Client {
	// 1. 初始化 option
	opts := &sms4goOptions{}

	// 2. 遍历用户选项
	for i := range options {
		options[i](opts)
	}

	// 3. 返回客户端
	return &smsClient{}
}

// GetSmsBlender 负载均衡获取 Blender
func GetSmsBlender() {

}

// GetSmsBlenderWithConfigId 根据 ConfigId 获取 Blender
func GetSmsBlenderWithConfigId(configId string) {

}
